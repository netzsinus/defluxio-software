#ifndef FREQ_CAPTURE_H
#define FREQ_CAPTURE_H 1

#include <stdint.h>

typedef enum {
  IDLE,
  PENDING,
  RUNNING
} mstate_t;

void freq_capture_init(void);
void freq_capture_start(void);
char freq_new_second(void);
mstate_t freq_get_state(void);
double freq_get_result(void);

#endif /* FREQ_CAPTURE_H */
