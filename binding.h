#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>
#include "Python.h"

// Convenience method for returning None
inline PyObject *PYNONE() {Py_INCREF(Py_None); return Py_None;}

// Receiver function in Go. Cant be called directly from python since we cannot
// return the *PyObject pointer without violating the cgo rules
extern char *_gosnake_invoke(PyObject *self, PyObject *args);

// If the result is nil attempts to print the error message
static void check_pyerr() {
    PyObject *err = PyErr_Occurred();

    if (err != NULL) {
        PyErr_PrintEx(0);
        Py_DECREF(err);
    }
}

// Called directly from Python. Invokes the go function _gosnake_invoke,
// runs the results through json.loads, then returns it to python
static PyObject* _gosnake_receive(PyObject *self, PyObject *args) {
    // Not sure that we need this, but leaving it here for now
    // PyGILState_STATE gil = PyGILState_Ensure();

    // route the call into Go
    char *result = _gosnake_invoke(self, args);
    PyObject *pyStringResult = PyString_FromString(result);

    // free the string
    free(result);

    // import json
    PyObject *json = PyImport_Import(PyString_FromString("json"));
    PyObject *loads = PyObject_GetAttr(json, PyString_FromString("loads"));

    // Build args tuple
    PyObject *tup = PyTuple_New(1);
    PyTuple_SET_ITEM(tup, 0, pyStringResult);

    // Make the call
    PyObject *_result = PyObject_CallObject(loads, tup);

    // Cleanup
    Py_DECREF(json);
    Py_DECREF(loads);
    Py_DECREF(tup);

    // Finish up and return
    check_pyerr();
    // PyGILState_Release(gil);
    return _result;
}

static PyMethodDef GosnakeMethods[] = {
    {"gocall", _gosnake_receive, METH_VARARGS, "doc string"},
    {NULL},
};

// Set up the python environment and inject gosnake as a virtual-python module
static void pyinit (int log) {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    Py_InitModule("gosnake", GosnakeMethods);
    PyEval_ReleaseThread(PyGILState_GetThisThreadState());

    if (log != 0) {
        fprintf(stdout, "gosnake: initialized python env\n");
    }
}