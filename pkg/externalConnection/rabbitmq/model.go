package rabbitmq

type MessageStats struct {
	Ack                 int  `json:"ack"`
	AckDetails          Rate `json:"ack_details"`
	Deliver             int  `json:"deliver"`
	DeliverDetails      Rate `json:"deliver_details"`
	DeliverGet          int  `json:"deliver_get"`
	DeliverGetDetails   Rate `json:"deliver_get_details"`
	DeliverNoAck        int  `json:"deliver_no_ack"`
	DeliverNoAckDetails Rate `json:"deliver_no_ack_details"`
	Get                 int  `json:"get"`
	GetDetails          Rate `json:"get_details"`
	GetEmpty            int  `json:"get_empty"`
	GetEmptyDetails     Rate `json:"get_empty_details"`
	GetNoAck            int  `json:"get_no_ack"`
	GetNoAckDetails     Rate `json:"get_no_ack_details"`
	Publish             int  `json:"publish"`
	PublishDetails      Rate `json:"publish_details"`
	Redeliver           int  `json:"redeliver"`
	RedeliverDetails    Rate `json:"redeliver_details"`
}

type BackingQueueStatus struct {
	AvgAckEgressRate  float64       `json:"avg_ack_egress_rate"`
	AvgAckIngressRate float64       `json:"avg_ack_ingress_rate"`
	AvgEgressRate     float64       `json:"avg_egress_rate"`
	AvgIngressRate    float64       `json:"avg_ingress_rate"`
	Delta             []interface{} `json:"delta"`
	Len               int           `json:"len"`
	Mode              string        `json:"mode"`
	NextDeliverSeqId  int           `json:"next_deliver_seq_id"`
	NextSeqId         int           `json:"next_seq_id"`
	NumPendingAcks    int           `json:"num_pending_acks"`
	NumUnconfirmed    int           `json:"num_unconfirmed"`
	Q1                int           `json:"q1"`
	Q2                int           `json:"q2"`
	Q3                int           `json:"q3"`
	Q4                int           `json:"q4"`
	TargetRamCount    string        `json:"target_ram_count"`
	Version           int           `json:"version"`
}

type GarbageCollection struct {
	FullsweepAfter  int `json:"fullsweep_after"`
	MaxHeapSize     int `json:"max_heap_size"`
	MinBinVheapSize int `json:"min_bin_vheap_size"`
	MinHeapSize     int `json:"min_heap_size"`
	MinorGcs        int `json:"minor_gcs"`
}

type Rate struct {
	Rate float64 `json:"rate"`
}

type Queue struct {
	Arguments                     map[string]interface{} `json:"arguments"`
	AutoDelete                    bool                   `json:"auto_delete"`
	BackingQueueStatus            BackingQueueStatus     `json:"backing_queue_status"`
	ConsumerCapacity              float64                `json:"consumer_capacity"`
	ConsumerUtilisation           float64                `json:"consumer_utilisation"`
	Consumers                     int                    `json:"consumers"`
	Durable                       bool                   `json:"durable"`
	EffectivePolicyDefinition     map[string]interface{} `json:"effective_policy_definition"`
	Exclusive                     bool                   `json:"exclusive"`
	ExclusiveConsumerTag          string                 `json:"exclusive_consumer_tag"`
	GarbageCollection             GarbageCollection      `json:"garbage_collection"`
	HeadMessageTimestamp          string                 `json:"head_message_timestamp"`
	IdleSince                     string                 `json:"idle_since"`
	Memory                        int                    `json:"memory"`
	MessageBytes                  int                    `json:"message_bytes"`
	MessageBytesPagedOut          int                    `json:"message_bytes_paged_out"`
	MessageBytesPersistent        int                    `json:"message_bytes_persistent"`
	MessageBytesRam               int                    `json:"message_bytes_ram"`
	MessageBytesReady             int                    `json:"message_bytes_ready"`
	MessageBytesUnacknowledged    int                    `json:"message_bytes_unacknowledged"`
	MessageStats                  MessageStats           `json:"message_stats"`
	Messages                      int                    `json:"messages"`
	MessagesDetails               Rate                   `json:"messages_details"`
	MessagesPagedOut              int                    `json:"messages_paged_out"`
	MessagesPersistent            int                    `json:"messages_persistent"`
	MessagesRam                   int                    `json:"messages_ram"`
	MessagesReady                 int                    `json:"messages_ready"`
	MessagesReadyDetails          Rate                   `json:"messages_ready_details"`
	MessagesReadyRam              int                    `json:"messages_ready_ram"`
	MessagesUnacknowledged        int                    `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails Rate                   `json:"messages_unacknowledged_details"`
	MessagesUnacknowledgedRam     int                    `json:"messages_unacknowledged_ram"`
	Name                          string                 `json:"name"`
	Node                          string                 `json:"node"`
	OperatorPolicy                string                 `json:"operator_policy"`
	Policy                        string                 `json:"policy"`
	RecoverableSlaves             string                 `json:"recoverable_slaves"`
	Reductions                    int                    `json:"reductions"`
	ReductionsDetails             Rate                   `json:"reductions_details"`
	SingleActiveConsumerTag       string                 `json:"single_active_consumer_tag"`
	State                         string                 `json:"state"`
	Type                          string                 `json:"type"`
	Vhost                         string                 `json:"vhost"`
}
