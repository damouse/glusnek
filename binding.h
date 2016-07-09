#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

// Convenience method for returning None
inline PyObject *PyNone() {Py_INCREF(Py_None); return Py_None;}

// A call from python to go. This is implemented below in go
extern PyObject* pyInvocation(PyObject *self, PyObject *args);

// Part of the threading implmentation
extern void createThreadCallback();
static void sig_func(int sig);


static PyObject* Foo_doSomething(PyObject *self, PyObject *args){
    printf("DoSomething called!\n");
    Py_INCREF(Py_None);
    return Py_None;
}


static PyMethodDef ModuleMethods[] = {
    {"summat", Foo_doSomething, METH_VARARGS, "doc string"},
    {"gocall", pyInvocation, METH_VARARGS, "doc string"},
    {NULL},
};

// Set up the python environment and inject gosnake as a virtual-python module
static void initialize_python (int log) {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    Py_InitModule("gosnake", ModuleMethods);
    PyEval_ReleaseThread(PyGILState_GetThisThreadState());

    if (log != 0) {
        fprintf(stdout, "gosnake: initialized python env\n");
    }
}

// Threading implementation
static void createThread(pthread_t* pid) {
    pthread_create(pid, NULL, (void*)createThreadCallback, NULL);
}

static void sig_func(int sig) {
    fprintf(stdout, "gosnake: handling exit signal\n");

    // PyObject *ptype, *pvalue, *ptraceback;
    // PyErr_Fetch(&ptype, &pvalue, &ptraceback);
    // pvalue contains error message
    // ptraceback contains stack snapshot and many other information

    // Get error message
    // char *pStrErrorMessage = PyString_AsString(pvalue);
    // fprintf(stdout, "pStrErrorMessage\n");

    signal(SIGSEGV, sig_func);
    pthread_exit(NULL);
}

static void register_sig_handler() {
    signal(SIGSEGV,sig_func);
}

