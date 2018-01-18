//   ,    ,  ,     __   _,  ___,
//   |    \_/     '|_) / \,' |
//  '|__ , /`     _|_)'\_/   |
//     '(_/      '     '     '
//
// a little robotic platform build for my son.
// It is meant to teach a bit of everything electronics related to a 6 year old in a playfull way.
// Most hardware choices are made out of accesibility considerations ( whatever we have lying around will do)
// As base we use a tank chassis driven by 2 DC motors. For motor control a DOIT 2 motor + 16 servo shield is used.
// might be changed later

//Motor A
const int motorVelA  = 9;
const int motorDirPinA  = 8;
//Motor B
const int motorDirPinB  = 7;
const int motorVelB  = 6;

const int deadbandUpper =  70;
const int deadbandLower = -70; 

boolean motorAfwd = true ;
boolean motorBfwd = true ;

// for initial control we use a X4S-SR SBUS receiver as it is convenient 
// to have long distance radio control 

#include <FUTABA_SBUS.h>

FUTABA_SBUS sBus;

void setup(){
  Serial.begin(115200);
  
  sBus.begin();

  pinMode(motorVelA, OUTPUT);
  pinMode(motorDirPinA, OUTPUT);
  pinMode(motorDirPinB, OUTPUT);
  pinMode(motorVelB, OUTPUT);

}

void loop(){
  //Serial.println("Loop");
  delay(50);
  sBus.FeedLine();
  if (sBus.toChannels == 1){
    sBus.UpdateServos();
    sBus.UpdateChannels();
    sBus.toChannels = 0;

    boolean motorAfwd = true ;
    boolean motorBfwd = true ;
    
    int motorX = (sBus.channels[1]);
    int motorY = (sBus.channels[2]);
    int motorA = map(motorX,172,1811,-255,255);
    int motorB = map(motorY,187,1811,-255,255);
    
    int motorspdA = (motorB + motorA);
    motorspdA = constrain(motorspdA,-255,255);
    
    int motorspdB = (motorB - motorA);
    motorspdB = constrain(motorspdB,-255,255);

    if (motorspdA < deadbandLower){ 
      //
      motorAfwd = false ;
      motorspdA = map(motorspdA,-255,0,255,0); 
    }
    if (motorspdA < deadbandUpper){ motorspdA = 0; }

    if (motorspdB < deadbandLower){ 
      motorBfwd = false ;
      motorspdB = map(motorspdB,-255,0,255,0);
      }
    if (motorspdB < deadbandUpper){ motorspdB = 0; }
    
    motorControl(motorspdA,motorAfwd,motorspdB,motorBfwd);

// we Print stuff to the serial port 

    Serial.print(motorspdA);
    Serial.print("\t");
    Serial.print(motorspdB);
    Serial.print("\t");
    Serial.print(motorAfwd);
    Serial.print("\t");
    Serial.print(motorBfwd);
    Serial.print("\t");
    Serial.print(sBus.channels[1]);
    Serial.print("\t");
    Serial.print(sBus.channels[2]);
    Serial.print("\t");
    Serial.print(sBus.channels[3]);
    Serial.print("\t");
    Serial.println(sBus.channels[4]);
    }
    
}

void motorControl(int velA, boolean dirA, int velB, boolean dirB){

    analogWrite(motorVelA, velA);
    digitalWrite(motorDirPinA, dirA);
    digitalWrite(motorDirPinB, !dirB);
    analogWrite(motorVelB, velB);
  
}
