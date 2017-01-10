import RPi.GPIO as GPIO, sys, threading, time, os, subprocess
sonar = 38
GPIO.setwarnings(False)
GPIO.setmode(GPIO.BOARD)
GPIO.setup(sonar, GPIO.OUT)
GPIO.output(sonar, True)
time.sleep(0.00001)
GPIO.output(sonar, False)
start = time.time()
count=time.time()
GPIO.setup(sonar,GPIO.IN)
while GPIO.input(sonar)==0 and time.time()-count<0.1:
    start = time.time()
count=time.time()
stop=count
while GPIO.input(sonar)==1 and time.time()-count<0.1:
    stop = time.time()
elapsed = stop-start
distance = elapsed * 34000
distance = distance / 2
print(distance)
GPIO.cleanup()
