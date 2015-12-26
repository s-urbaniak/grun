#include <glib.h>

extern int run();

static void idle_add() {
  g_idle_add(run, NULL);
}
