//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	guestSlice := guests[:]

	uniqueLoads := make(map[int]*Load)
	for _, g := range guestSlice {
		if uniqueLoads[g.CheckInDate] != nil {
			uniqueLoads[g.CheckInDate].GuestCount++
		} else {
			uniqueLoads[g.CheckInDate] = &Load{
				StartDate:  g.CheckInDate,
				GuestCount: 1,
			}
		}
		if uniqueLoads[g.CheckOutDate] != nil {
			uniqueLoads[g.CheckOutDate].GuestCount--
		} else {
			uniqueLoads[g.CheckOutDate] = &Load{
				StartDate:  g.CheckOutDate,
				GuestCount: -1,
			}
		}
	}

	loads := make([]*Load, len(uniqueLoads))
	i := 0
	for _, v := range uniqueLoads {
		loads[i] = v
		i++
	}

	sort.Slice(loads, func(i, j int) bool {
		return loads[i].StartDate < loads[j].StartDate
	})

	currectCount := 0
	for i := range loads {
		currectCount += loads[i].GuestCount
		loads[i].GuestCount = currectCount
	}

	loadsUnique := make([]Load, 0)
	lastLoadCount := -1
	for _, l := range loads {
		if lastLoadCount != l.GuestCount {
			lastLoadCount = l.GuestCount
			loadsUnique = append(loadsUnique, *l)
		}
	}

	return loadsUnique
}
