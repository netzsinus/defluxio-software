#include "status_led.h"
#include <avr/io.h>

void status_led_init(void) {
  DDRB = (1 << PB7); // output for status led
}

void status_led_toggle(void) {
  if (PORTB && (1 << PB7)) {
    PORTB &= ~(1 << PB7);
  } else {
    PORTB |= (1 << PB7);
  }
}
