# Python 3.3.3 and 2.7.6
# python fo.py

from threading import Lock
from threading import Thread


i = 0
lock = Lock()

def incrementingFunction():
    global i
    lock.acquired()
    for j in range(0, 1000000):
        i = i + 1
    lock.release()


def decrementingFunction():
    global i
    lock.acquired()
    for j in range(0, 1000000):
        i = i - 1
    lock.release()



def main():
    global i

    incrementing = Thread(target=incrementingFunction, args=(), )
    decrementing = Thread(target=decrementingFunction, args=(), )

    incrementing.start()
    incrementing.join()

    decrementing.start()
    decrementing.join()

    print("The magic number is %d" % (i))



main()
