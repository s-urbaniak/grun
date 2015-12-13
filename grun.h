#include <glib.h>

extern void runFunc();

static gboolean run(gpointer user_data) {
  runFunc();
  return TRUE;
}

static void add_runner() {
  g_idle_add(run, NULL);
}
