// Test Object
//
// # Overview
// A Test object is a test of an object.
//
// List of features:
//
// * Feature 1.
// * Feature 2.
// * Feature 3.
// * Feature 4.
//
// # Usage
// Create a new text by calling ```test_new()```

// This not part of the top of file comment.

#ifndef TEST_H_
#define TEST_H_

#include <stdbool.h>
#include <stdint.h>

// Test states.
typedef enum {
    // Initial state.
    TEST_STATE_INIT = 0,
    TEST_STATE_IDLE = 1,
    TEST_STATE_RUNNING,

    // Test passed.
    TEST_STATE_PASSED = 1000,
    TEST_STATE_FAILED,  // Test failed.
} test_state_t;

// A test structure
//
// * List item 1.
//   * List item 1.1.
//   * List item 1.3.
// * List item 2
//
// Represents a test.
typedef struct test {
  int priority;  // Priority.
  float weight; // Weight.

  // Next object in list.  Null if last object.
  struct test *next;

  // TODO(konkers): support embedded unions.
  // TODO(konkers): support anonymous unions.
  // TODO(konkers): support embedded structs.
} test_t;

// An anonymous typedef struct.
typedef struct {
  int type;  // Object type.
} object_t;

// Can be an int or float.
typedef union {
  int int_val;
  float float_val;
} value_t;

// A function-like macro.
#define add(a, b) ((a) + (b))

// A value constant macro.
#define TIMER_CR_EN (1 << 5)

// TODO(konkers) recognize consts.
const uint32_t TEST_TIMEOUT = 1000;  // Timeout in [ms].

// {private}
//
// TODO(konkers): Implement private tagging;
extern bool g_test_enabled;

// A function prototype.
test_t *test_new(void);

// Sets the test priority.
//
// Args:
//  t         Test to operate on.
//  priority  Priority of the test.
//
// Sets the priority of <t> to <priority>.
void test_set_priority(test_t *t, int priority);

// A static inline function.
//
// Now with more comment lines!
static inline float test_weight(const test_t *t) {
    return t->weight;
}

#endif  // TEST_H_
