package helpers

import (
	"encoding/json"
	"fmt"
)

func GroupByDateTime(messages []*Line) []*Line {
	var groupedMessages []*Line
	groupedMap := make(map[string][]*Line)

	for _, message := range messages {
		if message.Time == "" {
			continue
		}

		_, exists := groupedMap[message.Time]
		if exists {
			groupedMap[message.Time] = append(groupedMap[message.Time], &Line{
				// time cleanned to group
				Name:    message.Name,
				Message: message.Message,
			})
		} else {
			groupedMap[message.Time] = []*Line{{
				// time cleanned to group
				Name:    message.Name,
				Message: message.Message,
			}}
		}
	}

	// sort by datetime
	// TODO

	for timeRange, messages := range groupedMap {

		fmt.Printf("[%s] %d messages\n", timeRange, len(messages))
		sampledmessages := make([]*Line, 4)
		for i := range len(messages) {
			if i < 4 {
				sampledmessages[i] = messages[i]
			} else {
				break
			}
		}
		b, _ := json.Marshal(sampledmessages)
		fmt.Printf("sampledmessages: %s\n", string(b))

		groupedMessages = append(groupedMessages, &Line{
			Time: timeRange,
		})

		groupedMessages = append(groupedMessages, messages...)
	}

	return groupedMessages
}
