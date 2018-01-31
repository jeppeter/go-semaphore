package semaphore

/*
#include <fcntl.h>
#include <sys/stat.h>
#include <semaphore.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include <stdio.h>
#include <time.h>

int get_errno() {
	int ret = errno;
	if (ret > 0) {
		ret = -ret;
	}
	if (ret == 0) {
		ret = -1;
	}
    return ret;
}

int set_errno(int ret) {
	if (ret >= 0) {
		errno = ret;
	} else {
		errno = -ret;
	}
	return;
}

void* _sem_open(char* name, int value) {
    sem_t* sem = sem_open((const char*)name, O_CREAT, 0644, value);
    if (sem == SEM_FAILED) {
    	return NULL;
    }
    return sem;
}
int _sem_close(void* sem) {
	int ret=0;
	if (sem != NULL) {
		ret=sem_close((sem_t*)sem);
		if (ret != 0) {
			ret = get_errno();
		}
	}
    return ret;
}
int _sem_wait(void* sem) {
	int ret = 1;
	if (sem != NULL) {
		ret = sem_wait((sem_t*)sem);
		if (ret != 0) {
			ret = get_errno();
		}
	}
    return ret;
}
int _sem_trywait(void* sem) {
	int ret = -1;
	if (sem != NULL) {
		ret = sem_trywait((sem_t*)sem);
		if (ret != 0) {
			ret = get_errno();
		}
	}
    return ret;
}
int _sem_post(void* sem) {
	int ret = -1;
	if (sem != NULL) {
		ret = sem_post((sem_t*)sem);
		if (ret != 0) {
			ret = get_errno();
		}
	}
    return ret;
}

#define MAX_LONGLONG 0xffffffffffffffffULL

unsigned long long get_tickcount() {
	int ret;
	unsigned long long lret=MAX_LONGLONG;
	struct timespec tp;
	ret = clock_gettime(CLOCK_MONOTONIC,&tp);
	if (ret >= 0) {
		lret = tp.tv_sec * 1000;
		lret += (tp.tv_nsec / 1000000);
	}
	return lret;
}

int expire_time(unsigned long long sticks ,unsigned long long cticks,int mills){
	int ret=  0;
	int exmills;
	if ((MAX_LONGLONG - sticks) < (unsigned long long)mills) {
		if (cticks < sticks) {
			exmills = (int)cticks;
			exmills += (int)(MAX_LONGLONG - sticks);
			if (exmills > mills) {
				ret = 1;
			}
		}

	} else {
		if ((cticks - sticks) > mills) {
			ret = 1;
		}
	}
	return ret;
}

int _sem_waittimed(void* sem,int mills) {
	unsigned long long sticks,cticks;
	int ret;
	fd_set rset;
	int maxfd;
	struct timeval tmval;

	sticks = get_tickcount();
	while(1) {
		ret = _sem_trywait(sem);
		if (ret >= 0) {
			return 0;
		}
		cticks = get_tickcount();
		ret = expire_time(sticks,cticks,mills);
		if (ret > 0) {
			set_errno(-ETIMEDOUT);
			return -1;
		}
		maxfd = fileno(stdin) + 1;
		FD_ZERO(&rset);
		FD_SET(fileno(stdin),&rset);
		tmval.tv_sec = 0;
		tmval.tv_usec = 1000;
		ret = select(maxfd,&rset,NULL,NULL,&tmval);
	}
	return -1;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Semaphore struct {
	Name string
	hdl  uintptr
}

func NewSemaphore(name string, cnt int) (*Semaphore, error) {
	var p *Semaphore
	var err error
	var fname string
	fname = fmt.Sprintf("/%s", name)
	cName := C.CString(fname)
	p = &Semaphore{}
	p.Name = name
	p.hdl = uintptr(C._sem_open(cName, C.int(cnt)))
	if p.hdl == 0 {
		err = fmt.Errorf("can not create [%s] error[%d]", C.get_errno())
		return nil, err
	}
	return p, nil
}

func (psema *Semaphore) Wait(mills int) error {
	if psema.hdl == uintptr(0) {
		return fmt.Errorf("not valid hdl")
	}

	if mills < 0 {
		errint := C._sem_wait(unsafe.Pointer(psema.hdl))
		if errint != 0 {
			return fmt.Errorf("can not wait on [%s] error[%v]", psema.Name, C.get_errno())
		}
	} else {
		errint := C._sem_waittimed(unsafe.Pointer(psema.hdl), C.int(mills))
		if errint != 0 {
			return fmt.Errorf("can not wait on [%s] time error[%v]", psema.Name, C.get_errno())
		}
	}
	return nil
}

func (psema *Semaphore) Release() error {
	if psema.hdl == uintptr(0) {
		return fmt.Errorf("not valid hdl")
	}
	errint := C._sem_post(unsafe.Pointer(psema.hdl))
	if errint != 0 {
		return fmt.Errorf("post [%s] error[%v]", psema.Name, C.get_errno())
	}
	return nil
}

func (psema *Semaphore) Close() {
	if psema.hdl != uintptr(0) {
		C._sem_close(unsafe.Pointer(psema.hdl))
		psema.hdl = uintptr(0)
	}
	return
}
