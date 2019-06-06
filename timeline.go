package TimeLine

import (
	"fmt"
	"sort"
)

type (
	// EventTime время события
	EventTime struct {
		// Begin - начало события, End - конец события
		Begin, End int
	}
	// TimeLine - временная линия
	TimeLine struct {
		// EventTimes список времен событий
		EventTimes []EventTime
		// Day - начало и конец временной линии (например, рабочий день)
		Day EventTime
	}
)

// String - переопределение для вывода времени события в формате hh:mm – hh:mm
func (e EventTime) String() string {
	startH, startM := tL2time(e.Begin)
	endH, endM := tL2time(e.End)
	return fmt.Sprintf("%.2d:%.2d – %.2d:%.2d", startH, startM, endH, endM)
}

// CreateTL - конструктор создания временной линии
// beginHour, beginMinute, endHour, endMinute - время начала и конца временной линии в часах и минутах
func CreateTL(beginHour, beginMinute, endHour, endMinute int) (TimeLine, error) {
	if beginHour < 0 || beginHour > 23 {
		beginHour = 0
	}
	if beginMinute < 0 || beginMinute > 59 {
		beginMinute = 0
	}
	if endHour < 0 || endHour > 23 {
		endHour = 23
	}
	if endMinute < 0 || endMinute > 59 {
		endMinute = 59
	}
	beginTL := time2tL(beginHour, beginMinute)
	endTL := time2tL(endHour, endMinute)
	if beginTL > endTL {
		return TimeLine{}, fmt.Errorf("the beginning of the period is later than the ending")
	}
	tl := TimeLine{Day: EventTime{Begin: beginTL, End: endTL}}
	return tl, nil
}

// Add - добавление периода события во временную линию
// hoursBegin, minutesBegin, hoursEnd, minutesEnd - время начала и конца события в часах и минутах
// doNotMatter - добавлять ли событие, если его время пересекается со временем добавленного ранее события
// возвращает ошибку, если время пересекается
func (tl *TimeLine) Add(hoursBegin, minutesBegin, hoursEnd, minutesEnd int, doNotMatter bool) error {
	return tl.addEvent(time2tL(hoursBegin, minutesBegin), time2tL(hoursEnd, minutesEnd), doNotMatter)
}

func (tl *TimeLine) addEvent(begin, end int, doNotMatter bool) (err error) {
	for _, event := range tl.EventTimes {
		if (begin > event.Begin && begin < event.End) || (end > event.Begin && end < event.End) {
			err = fmt.Errorf("event intersects with other events")
			if doNotMatter {
				break
			} else {
				return
			}
		}
	}
	(*tl).EventTimes = append((*tl).EventTimes, EventTime{Begin: begin, End: end})
	tl.sort()
	return
}

// AddDuration - добавляет время события в свободное "окно"
// duration - длительность события в минутах
// first - признак добавления (true - добавить в первое подходящее "окно", false - добавить в минимальное подходящее "окно"
func (tl *TimeLine) AddDuration(duration int, first bool) (EventTime, error) {
	var (
		err        error
		begin, end int
	)
	events := tl.GetEmpty()
	if len(events) == 0 {
		return EventTime{}, fmt.Errorf("no free period")
	}
	var (
		index        = -1
		min          = 0
		tempDuration int
		flag         = true
	)
	for i := range events {
		tempDuration = events[i].End - events[i].Begin
		if (tempDuration > duration) && (flag || tempDuration < min) {
			if first {
				begin, end = events[index].Begin, events[index].Begin+duration
				err = tl.addEvent(begin, end, false)
				return EventTime{begin, end}, err
			}
			index = i
			min = tempDuration
			flag = false
		}
	}
	if index == -1 {
		return EventTime{}, fmt.Errorf("no free period")
	}
	begin, end = events[index].Begin, events[index].Begin+duration
	err = tl.addEvent(begin, end, false)
	return EventTime{begin, end}, err
}

// GetEmpty() получить список свободных "окон" во временной линии
func (tl TimeLine) GetEmpty() []EventTime {
	events := make([]EventTime, 0, 20)
	var begin, end int
	if tl.EventTimes[0].Begin > tl.Day.Begin {
		begin = tl.Day.Begin
		end = tl.EventTimes[0].Begin
		events = append(events, EventTime{Begin: begin, End: end})
	}
	for i := 0; i < len(tl.EventTimes)-1; i++ {
		if tl.EventTimes[i].End < tl.EventTimes[i+1].Begin {
			begin = tl.EventTimes[i].End
			end = tl.EventTimes[i+1].Begin
			events = append(events, EventTime{Begin: begin, End: end})
		}
	}
	if tl.EventTimes[len(tl.EventTimes)-1].End < tl.Day.End {
		begin = tl.EventTimes[len(tl.EventTimes)-1].End
		end = tl.Day.End
		events = append(events, EventTime{Begin: begin, End: end})
	}
	return events
}

// сортировка событий по времени начала
func (tl *TimeLine) sort() {
	sort.Slice((*tl).EventTimes, func(i, j int) bool { return (*tl).EventTimes[i].Begin < (*tl).EventTimes[j].Begin })
}

// конвертация времени из часов и минут в смещение в минутах от начала суток
func time2tL(hour, minute int) int {
	return hour*60 + minute
}

// конвертация времени из смещения в минутах от начала суток в часы и минуты
func tL2time(tl int) (int, int) {
	absTL := tl
	return absTL / 60, absTL % 60
}