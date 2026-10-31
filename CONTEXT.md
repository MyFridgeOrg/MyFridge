# MyFridge

A camera mounted at a fridge that takes a photo whenever someone is in front of it.

## Language

**Fridge unit**:
The Raspberry Pi together with its camera and PIR motion sensor, installed at a fridge.
_Avoid_: Device, the Pi (when meaning the whole setup)

**Motion event**:
One continuous period of motion reported by the PIR sensor, starting when its signal goes high. Retriggers while the signal is still high belong to the same motion event.
_Avoid_: Trigger, detection, motion

**Snapshot**:
One photo taken by the fridge unit in response to a motion event.
_Avoid_: Capture, image, picture (as nouns)
