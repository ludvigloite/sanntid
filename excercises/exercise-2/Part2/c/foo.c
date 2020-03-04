//
//  main.c
//  Sanntid_oving1
//
//  Created by Ludvig Løite on 17/01/2020.
//  Copyright © 2020 Ludvig Løite. All rights reserved.
//

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t lock;

// Note the return type: void*
void* incrementingThreadFunction(){
    pthread_mutex_lock(&lock); //låser
    for (int j = 0; j < 1000000; j++) {
        // TODO: sync access to i
        i++;
    }
    pthread_mutex_unlock(&lock); //låser opp
    return NULL;
}

void* decrementingThreadFunction(){
    pthread_mutex_lock(&lock); //låser
    for (int j = 0; j < 1000000; j++) {
        // TODO: sync access to i
        i--;
    }
    pthread_mutex_unlock(&lock); //låser opp
    return NULL;
}


int main(){
    
    //her initialiserer jeg en mutex med navn lock.
    if (pthread_mutex_init(&lock, NULL) != 0){
        printf("\n mutex init failed");
        return 1;
    }
    
    pthread_t incrementingThread, decrementingThread;
    
    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL);
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);
    
    pthread_join(incrementingThread, NULL);
    pthread_join(decrementingThread, NULL);
    
    pthread_mutex_destroy(&lock);
    
    printf("The magic number is: %d\n", i);
    return 0;
}
