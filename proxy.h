#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>

inline PyObject *PyNone() {Py_INCREF(Py_None); return Py_None;}

extern PyObject* PythonInvocation(PyObject *self, PyObject *args);

static PyMethodDef ModuleMethods[] = {
    {"gocall", PythonInvocation, METH_VARARGS, "doc string"},
    {NULL},
};

static void initialize_python () {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    Py_InitModule("gosnake", ModuleMethods);
    PyEval_ReleaseThread(PyGILState_GetThisThreadState());
    fprintf(stdout, "gosnake: initialized python env\n");
}