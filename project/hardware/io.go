
//Gjør funksjonene fra io.c om til go kode. Dette er veldig enkelt.

import "C"


func io_init() int{
	return int(C.io_init())
}

func io_set_bit(channel int){
	C.io_set_bit(C.int(channel))
}

//TODO: fiks resten i liknende stil





int io_init(void);

void io_set_bit(int channel);
void io_clear_bit(int channel);

int io_read_bit(int channel);

int io_read_analog(int channel);
void io_write_analog(int channel, int value);