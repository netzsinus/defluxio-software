/*************************************************************************
Title:    example program for the Interrupt controlled UART library
Author:   Peter Fleury <pfleury@gmx.ch>   http://jump.to/fleury
File:     $Id: test_uart.c,v 1.4 2005/07/10 11:46:30 Peter Exp $
Software: AVR-GCC 3.3
Hardware: any AVR with built-in UART, tested on AT90S8515 at 4 Mhz

DESCRIPTION:
          This example shows how to use the UART library uart.c

*************************************************************************/
#include <stdlib.h>
#include <avr/io.h>
#include <avr/interrupt.h>
#include <avr/pgmspace.h>
#include <stdio.h>
#include <math.h>
#include <string.h>
#include "uart.h"

/* define CPU frequency in Mhz here if not defined in Makefile */
#ifndef F_CPU
#define F_CPU 16000000UL
#endif

/* 9600 baud */
#define UART_BAUD_RATE      9600      

#define SAMPLES 1000
#define t1 1 // 1 ms.
#define t2 1000 // 500 ms

#define Vref 5.0


long FP_ADC_CONV, FP_OFFSET; 

////// 20:12 fixed point macros
#define int2fix(a)   (((long)(a))<<12)    //Convert char to fix. a is a char
#define float2fix(a) ((long)((a)*4096.0)) //Convert float to fix. a is a float
#define fix2float(a) ((float)(a)/4096.0)  //Convert fix to float. a is an int
#define multfix(a,b) ((long)((((long)(a))*((long)(b)))>>12))

long fpVSquareSum;
long fpISquareSum;
long fpPowerSum;
float rPower, aPower;
float powerFactor;

float iScaleFactor, vScaleFactor, rpScaleFactor;

volatile unsigned int msecs;
///// Periods //////
long fpVLastValue;
unsigned long vPeriodSum;
int vLastZeroMsec;
unsigned int vPeriodCount;


///// Watt Hour Stuff ////
float kWattHrs;

volatile long time0, time1, time2;
unsigned int sample;

long threshold;
char shutdown;

int fpMult(int a, int b);



void task2(void) {
  time2 = t2;
  if (PORTB && (1 << PB7)) {
    PORTB &= ~(1 << PB7);
  } else {
    PORTB |= (1 << PB7);
  }
}

void task1(void) {
  char Ain0 = 0;
  long fpVin0 = 0; 
  long fpAin0 = 0;
  float vAvg = 0;
  int vPeriod = 0;


  // start conversion for channel 0
  ADMUX &= ~(1 << MUX0);
  ADCSRA |= (1<<ADSC);                
  // wait for conversion to finish
  while (ADCSRA & (1<<ADSC) ) {
  }
  // read Ain0;
  Ain0 = ADCW;
  
  // go to 20:12 fixed point
  fpAin0 = (long)Ain0 << 12;

  fpVin0 = multfix(fpAin0, FP_ADC_CONV) - FP_OFFSET;
//  fpVin = multfix(fpVin0, FP_V_CONV); // convert to actual voltage

  fpVSquareSum += multfix(fpVin0, fpVin0);

  ///////// FREQUENCY CALCULATION //////
  if (fpVLastValue < 0 && fpVin0 >= 0) {
    // found a zero crossing from neg. to pos.!
    vPeriod = msecs - vLastZeroMsec;

    if (vPeriod > 0) {
      vPeriodSum += vPeriod;
      vPeriodCount++;
    }

    vLastZeroMsec = msecs;

  }

  fpVLastValue = fpVin0;
  ///////////////////////////////////


  sample++;

  if (sample == SAMPLES) {
    char buffer[10];
    // this only happens once a second,
    // so floating point mults should be okay
    vAvg = (float) sqrt(vScaleFactor * fpVSquareSum);

    //if (power < 10) power = 0; // probably just noise

    //sprintf(t_buffer, "P%3.2f,V%3.2f,F%3.2f\n\r", rPower, vAvg, 1000.0 * (float)vPeriodCount / (float)vPeriodSum);
    uart_puts("M;");
    dtostrf(vAvg, 8, 4, buffer);
    uart_puts(buffer);
    uart_puts(";");
    dtostrf(1000.0 * (float)vPeriodCount / (float) vPeriodSum, 8, 4, buffer);
    uart_puts(buffer);
    uart_puts(";");
    dtostrf(fpVin0, 8, 4, buffer);
    uart_puts(buffer);
    uart_puts("\r\n");

    // reset for next sample
    fpVSquareSum = 0;
    fpISquareSum = 0;
    fpPowerSum = 0;
    vPeriodSum = 0;
    vPeriodCount = 0;
    sample = 0;
  }

  time1 = t1;
}

//interrupt [TIM0_COMP] void timer0_compare(void) {
ISR(TIMER0_COMPA_vect) {
  if (time1 > 0) --time1;
  if (time2 > 0) --time2;
  msecs++;
  if (msecs == 1000)
  {
    msecs = 0;
  }
}

void initialize(void) {
  /////////////// TIMER 0 ///////////////////////
  TIMSK0 |= (1 << OCIE0A); // turn on timer 0 compare match ISR
  OCR0A = 250; // sets compare to 250 ticks
  TCCR0A |= (1 << WGM01); // Mode 2/CTC
  TCCR0B |= (1 << CS01) | (1 << CS00); // prescaler 64

  msecs = 0;

  DDRB = (1 << PB7); // output for status led

  time1 = t1;
  time2 = t2;

  // calculate conversion factor
  vScaleFactor = 2002001.0 / (4096.0 * SAMPLES);
  //FP_ADC_CONV = float2fix(Vref / 256.0);
  FP_ADC_CONV = float2fix(0.512 / 256.0); // opto
  FP_OFFSET = float2fix(0.256);

  // internal Aref = Vcc, left adjust result
  ADMUX = (1<<REFS0) | (1 << ADLAR);
  //enable ADC and set prescaler to 1/128*16MHz=125,000
  //and clear interupt enable
  ADCSRA = ((1<<ADPS2) | (1<<ADPS1) | (1<<ADPS0));     // Frequenzvorteiler
  ADCSRA |= (1<<ADEN);                  // ADC aktivieren
 
  // Dummy readout to get the ADC going
  ADCSRA |= (1<<ADSC);                
  while (ADCSRA & (1<<ADSC) ) {
  }
  // Read ADCW once
  (void) ADCW;

  uart_init( UART_BAUD_SELECT(UART_BAUD_RATE,F_CPU) ); 
  sei();
}

int main(void) {
  initialize();
  uart_puts_P("I;defluxio startup complete\n\r");

  while(1) {
    if (time1 == 0) task1();
    if (time2 == 0) task2();
  }
}
