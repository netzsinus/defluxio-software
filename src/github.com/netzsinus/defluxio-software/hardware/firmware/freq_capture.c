#include "freq_capture.h"
#include "status_led.h"
#include <avr/io.h>
#include <avr/interrupt.h>

#ifndef TRUE
#define TRUE 0
#define FALSE 1
#endif

volatile uint32_t end_time = 0;
volatile uint32_t start_time = 0;
volatile uint16_t timer1_overflows = 0;
volatile uint16_t edges = 0;

// timekeeping.
volatile uint16_t milliseconds = 0; 
volatile unsigned char new_second = FALSE;

// maintain the state of a measurement
volatile mstate_t m_state = IDLE;

void freq_capture_init(void) {
  /** PD6 is the input capture pin. Set data direction to input. */
  DDRD &= ~(1 << PD6);
  /** No internal pullup - we have that covered in hardware. */
  PORTD &= ~(1 << PD6);
  /** We use Timer1 for the input capture pin. 
   * - No Prescaler.
   * - Activate input capture noise canceler. 
   * - Trigger on rising edge.
   * */
  TCCR1B |= ((1<<ICES1) | (1 << ICNC1) | (1<<CS10));
  //TCCR1B |= ((1<<ICES1) | (1<<CS10));
  // Enable Capture & Overflow interrupts for Timer1
  TIMSK1 |= ((1<<ICIE1) | (1<<TOIE1));

  // Use Timer2 for timekeeping. CTC with prescaler 256.
  TCCR2A |= (1<<WGM21);
  TCCR2B |= ((1<<CS22) | (1<<CS21)); 
  OCR2A = (F_CPU/256)/1000 - 1;
  // Enable compare time interrupt
  TIMSK2 |= (1<<OCIE2A);
}

void freq_capture_start(void) {
  new_second = FALSE;
  m_state = PENDING;
}

mstate_t freq_get_state(void) {
  return m_state;
}

double freq_get_result(void) {
  double result = timer1_overflows * 65536 + end_time - start_time;
  return (F_CPU) / result * edges;
}

/*****************************************************
 * Timer interrupts below
 ****************************************************/

ISR ( TIMER2_COMPA_vect ) {
  milliseconds++;
  if(milliseconds == 1000) {
    // Mark the beginning of a new second - triggers state transitions
    new_second = TRUE;
    milliseconds = 0;
  }
}

ISR( TIMER1_OVF_vect ) {
  timer1_overflows++;
}

ISR( TIMER1_CAPT_vect ) {
  //status_led_toggle();
  if(m_state != IDLE) {
    if ((m_state == PENDING) & (new_second == TRUE)) {
      // Start a new measurement - initialize state variables
      start_time = ICR1;
      end_time = start_time;
      timer1_overflows = 0;
      edges = 0;
      new_second = FALSE;
      m_state = RUNNING;
    } else if ((m_state == RUNNING) & (new_second == TRUE)) {
      // measurement is finished
      m_state = IDLE;
      new_second = FALSE;
      end_time = ICR1;
      edges++;
    } else if (m_state == RUNNING) {
      // measurement in progress, update state variables
      end_time = ICR1;
      edges++;
    } 
  }
}
