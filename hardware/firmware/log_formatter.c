#include <avr/pgmspace.h>
#include <stdlib.h>
#include "log_formatter.h"
#include "uart.h"

void log_info(char* msg) {
  uart_puts_P("I;");
  uart_puts(msg);
  uart_puts_P("\r\n");
}

void log_freq(double freq) {
  char buffer[16];
  uart_puts_P("F;");
  dtostrf(freq, 7, 5, buffer);
  uart_puts(buffer);
  uart_puts_P("\r\n");
}
