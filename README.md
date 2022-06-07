# ATM
---
atm is a application that allows a user to deposit or get money from an ATM.

## Introduction
---
Separate the ATM machine into 3 main parts:
- account: bank account management, should connect to remote server or database, but is implementation locally now
- cashier: cash management of the ATM machine, should be related with hardware, but is implementation is a virtual one now
- ui: user interface of the machine, should be graph ui, now is implementation in commandline and a mock for test

## hierarchy
---
- cmd
  - cli: main function of client
- internal
  - account: interface and data structure for account related
    - impl: implementation folder
      - local: local implementation of manager and delegate
  - atm: business layer, implementation most of the requirement, it is a composition of account manager, cashier and ui
  - cashier: interface and data structure for cashier related
    - impl: implementation folder
      - cashier20: implementation with cashier only allow $20
  - command: parse command line text into tokens and call function to process
  - ui: simply define a ui interface
    - impl: implementation folder
      - commandline: just print the output to standard output
      - mock: save output, used for test