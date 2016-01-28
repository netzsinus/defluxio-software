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
#include "freq_capture.h"
#include "status_led.h"
#include "log_formatter.h"

/* 9600 baud */
#define UART_BAUD_RATE      9600      



void initialize(void) {
  
  //// internal Aref = Vcc, left adjust result
  //ADMUX = (1<<REFS0) | (1 << ADLAR);
  ////enable ADC and set prescaler to 1/128*16MHz=125,000
  ////and clear interupt enable
  //ADCSRA = ((1<<ADPS2) | (1<<ADPS1) | (1<<ADPS0));     // Frequenzvorteiler
  //ADCSRA |= (1<<ADEN);                  // ADC aktivieren

  //// Dummy readout to get the ADC going
  //ADCSRA |= (1<<ADSC);                
  //while (ADCSRA & (1<<ADSC) ) {
  //}
  //// Read ADCW once
  //(void) ADCW;

  freq_capture_init();
  status_led_init();
  uart_init( UART_BAUD_SELECT(UART_BAUD_RATE,F_CPU) ); 
  sei();
}

int main(void) {
  initialize();
  log_info("defluxio startup complete");
  freq_capture_start();
  while(1) {
    if (freq_get_state() == IDLE) {
      status_led_toggle();
      double current_freq = freq_get_result();
      log_freq(current_freq);
      freq_capture_start();
    } 
  }
}
