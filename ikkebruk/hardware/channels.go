
//TODO: fiks alle som PORT_4, OBSTRUCTION osv.

const{
//in port 4
	PORT_4			=	3
	OBSTRUCTION		=	(0x300+23)
	STOP 			=	(0x300+22)
	BUTTON_COMMAND1	=	(0x300+21)
	BUTTON_COMMAND2         (0x300+20)
	BUTTON_COMMAND3         (0x300+19)
	BUTTON_COMMAND4         (0x300+18)
	BUTTON_UP1              (0x300+17)
	BUTTON_UP2              (0x300+16)
}

const{
	//in port 1
	PORT_1        2
	BUTTON_DOWN2            (0x200+0)
	BUTTON_UP3              (0x200+1)
	#define BUTTON_DOWN3            (0x200+2)
	#define BUTTON_DOWN4            (0x200+3)
	#define SENSOR_FLOOR1           (0x200+4)
	#define SENSOR_FLOOR2           (0x200+5)
	#define SENSOR_FLOOR3           (0x200+6)
	#define SENSOR_FLOOR4           (0x200+7)
}

//out port 3
#define PORT_3_SUBDEVICE        3
#define PORT_3_CHANNEL_OFFSET   8
#define PORT_3_DIRECTION        COMEDI_OUTPUT
#define MOTORDIR                (0x300+15)
#define LIGHT_STOP              (0x300+14)
#define LIGHT_COMMAND1          (0x300+13)
#define LIGHT_COMMAND2          (0x300+12)
#define LIGHT_COMMAND3          (0x300+11)
#define LIGHT_COMMAND4          (0x300+10)
#define LIGHT_UP1               (0x300+9)
#define LIGHT_UP2               (0x300+8)

//out port 2
#define PORT_2_SUBDEVICE        3
#define PORT_2_CHANNEL_OFFSET   0
#define PORT_2_DIRECTION        COMEDI_OUTPUT
#define LIGHT_DOWN2             (0x300+7)
#define LIGHT_UP3               (0x300+6)
#define LIGHT_DOWN3             (0x300+5)
#define LIGHT_DOWN4             (0x300+4)
#define LIGHT_DOOR_OPEN         (0x300+3)
#define LIGHT_FLOOR_IND2        (0x300+1)
#define LIGHT_FLOOR_IND1        (0x300+0)

//out port 0
#define MOTOR                   (0x100+0)

//non-existing ports (for alignment)
#define BUTTON_DOWN1            -1
#define BUTTON_UP4              -1
#define LIGHT_DOWN1             -1
#define LIGHT_UP4               -1