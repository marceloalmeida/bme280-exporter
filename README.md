# BME280 on RPi zero

# RPi zero pinout

* [Pin Numbering - Raspberry Pi Zero W](https://pi4j.com/1.2/pins/model-zerow-rev1.html)

## Find i2c bus

```sh
# i2cdetect -l
i2c-1	i2c       	bcm2835 I2C adapter             	I2C adapter
#
#   ^
#   |
# this will be the i2c bus
```

## Find i2c address

```sh
# i2cdetect -y 1
     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:          -- -- -- -- -- -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
70: -- -- -- -- -- -- -- 77
#
#                         ^
#                         |
#                   i2c address
```
