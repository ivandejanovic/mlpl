CC = gcc -Wall -m64 -std=c11
DCC = gcc -g3 -Wall -m64 -std=c11

SOURCE = src/
BUILD = build/

CODE = $(SOURCE)*.c
TARGET = $(BUILD)mlpl

OBJ = *.o
BAK = *~

clean:
	rm -f $(TARGET) $(OBJ) $(BAK) $(SOURCE)$(OBJ) $(SOURCE)$(BAK) $(BUILD)$(OBJ) $(BUILD)$(BAK)

release:
	$(CC) -o $(TARGET) $(CODE)

debug:
	$(DCC) -o $(TARGET) $(CODE)

all: clean debug