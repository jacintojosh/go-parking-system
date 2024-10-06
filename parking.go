package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type ParkingLotService interface {
	CreateParkingLot(slots string) (pL ParkingLot, err error)
	Park(plate string, color string) (err error)
	Status() (pL ParkingLot)
	RemoveCar(slot string) (err error)
	SearchColorPlateNumber(color string) (err error)
	SearchColorSlotNumber(color string) (err error)
	SearchPlateNumberSlot(plate string) (err error)
}

type ParkingLot []ParkingSlot

type ParkingSlot struct {
	Slot        int
	PlateNumber string
	Color       string
}

const CREATE_PARKING_LOT_COMMAND = "create_parking_lot"
const PARKING_LOT_STATUS_COMMAND = "status"
const PARK_COMMAND = "park"
const REMOVE_CAR_COMMAND = "leave"
const SEARCH_COLOR_FOR_PLATES_COMMAND = "plate_numbers_for_cars_with_colour"
const SEARCH_COLOR_FOR_SLOT_COMMAND = "slot_numbers_for_cars_with_colour"
const SEARCH_PLATE_FOR_SLOT_COMMAND = "slot_number_for_registration_number"

var LIST_OF_COMMANDS = []string{CREATE_PARKING_LOT_COMMAND, REMOVE_CAR_COMMAND, PARKING_LOT_STATUS_COMMAND, PARK_COMMAND, SEARCH_COLOR_FOR_PLATES_COMMAND, SEARCH_COLOR_FOR_SLOT_COMMAND, SEARCH_PLATE_FOR_SLOT_COMMAND}

func (p ParkingLot) CreateParkingLot(slots string) (pL ParkingLot, err error) {
	slotsI, err := strconv.Atoi(slots)
	if err != nil {
		fmt.Printf("invalid amount of parking slots \"%s\", please type in a valid number\n", slots)
		return p, err
	}
	p = make([]ParkingSlot, slotsI)
	for i := range p {
		p[i].Slot = i + 1
	}
	if p == nil {
		err = errors.New("failed in making parking lot")
	}
	return p, err
}

func (p ParkingLot) Park(plate string, color string) (err error) {
	// error checking for valid plate and valid color omitted
	for i := range p {
		if p[i].PlateNumber == "" && p[i].Color == "" {
			p[i].PlateNumber = plate
			p[i].Color = color
			fmt.Printf("Allocated Slot number %d\n", i+1)
			return nil
		}
	}
	if p != nil {
		fmt.Println("Sorry, parking lot is full")
	} else {
		err = errors.New("parking lot does no exist yet, invalid command")
	}
	return err
}

func (p ParkingLot) Status() (err error) {
	if p == nil {
		fmt.Println("parking lot does not exist yet, please create a parking lot first")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Slot No.\t|Plate Number\t|Colour")
	for i := range p {
		if p[i].PlateNumber != "" && p[i].Color != "" {
			fmt.Fprintf(w, "%d\t|%s\t|%s\n", i+1, p[i].PlateNumber, p[i].Color)
		}
	}
	w.Flush()
	return nil
}

func (p ParkingLot) RemoveCar(slot string) (err error) {
	slotI, err := strconv.Atoi(slot)
	if err != nil {
		fmt.Printf("invalid parking slot number \"%s\", please type in a valid number\n", slot)
		return err
	}
	if slotI > len(p) {
		fmt.Printf("parking slot does not exist")
		return nil
	}
	for i := range p {
		if p[i].Slot == slotI {
			p[i].PlateNumber = ""
			p[i].Color = ""
			fmt.Printf("Slot number %s, is free\n", slot)
			return nil
		}
	}
	fmt.Println("Not Found")
	return nil
}

func (p ParkingLot) SearchColorPlateNumber(color string) (err error) {
	var plateNumbers []string
	for i := range p {
		if p[i].Color == color {
			plateNumbers = append(plateNumbers, p[i].PlateNumber)
		}
	}
	if len(plateNumbers) > 0 {
		fmt.Println(strings.Join(plateNumbers, ", "))
	} else {
		fmt.Println("Not found")
	}
	return nil
}

func (p ParkingLot) SearchColorSlotNumber(color string) (err error) {
	var slotNumbers []int
	for i := range p {
		if p[i].Color == color {
			slotNumbers = append(slotNumbers, i+1)
		}
	}
	if len(slotNumbers) > 0 {
		fmt.Println(strings.Join(arrIntToString(slotNumbers), ", "))
	} else {
		fmt.Println("Not found")
	}
	return nil
}

func (p ParkingLot) SearchPlateNumberSlot(plate string) (err error) {
	for i := range p {
		if p[i].PlateNumber == plate {
			fmt.Println(strconv.Itoa(i + 1))
			return nil
		}
	}
	fmt.Println("Not found")
	return nil
}

func main() {
	var pL ParkingLot
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command, arg1, arg2 := parse(scanner.Text())
		switch command {
		case CREATE_PARKING_LOT_COMMAND:
			pL, err = pL.CreateParkingLot(arg1)
		case PARK_COMMAND:
			err = pL.Park(arg1, arg2)
		case PARKING_LOT_STATUS_COMMAND:
			err = pL.Status()
		case REMOVE_CAR_COMMAND:
			err = pL.RemoveCar(arg1)
		case SEARCH_COLOR_FOR_PLATES_COMMAND:
			err = pL.SearchColorPlateNumber(arg1)
		case SEARCH_COLOR_FOR_SLOT_COMMAND:
			err = pL.SearchColorSlotNumber(arg1)
		case SEARCH_PLATE_FOR_SLOT_COMMAND:
			err = pL.SearchPlateNumberSlot(arg1)
		case "help":
			fmt.Println("Commands are", strings.Join(LIST_OF_COMMANDS, ", "))
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Command not recognized, type help for list of commands")
		}
		if err != nil {
			handleError(err)
		}
	}
}

func arrIntToString(a []int) (b []string) {
	b = make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return b
}

func parse(s string) (main string, arg1 string, arg2 string) {
	s = strings.TrimSpace(s)
	split := strings.Split(s, " ")
	main = split[0]
	if len(split) > 2 {
		arg1, arg2 = split[1], split[2]
	} else if len(split) == 2 {
		arg1 = split[1]
	}
	return main, arg1, arg2
}

func handleError(e error) {
	fmt.Println(e)
	fmt.Println("fatal error detected, stopping program")
	os.Exit(3)
}
