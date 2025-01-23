package kafka

import (
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/parsers"
)

const DefaultAuth = "admin"

type KafkaSource struct {
	Connection  *KafkaConnectionOptions     `json:"connection"`
	Auth        *KafkaAuth                  `json:"auth"`
	Topic       string                      `json:"topic"`
	GroupTopics []string                    `json:"group_topics"`
	Transformer *model.DataTransformOptions `json:"transformer"`

	BufferSize model.BytesSize `json:"buffer_size"` // it's not some real buffer size - see comments to waitLimits() method in kafka-source

	SecurityGroupIDs []string `json:"security_group_ids"`

	ParserConfig        map[string]interface{} `json:"parser_config"`
	IsHomo              bool                   `json:"is_homo"`               // enabled kafka mirror protocol which can work only with kafka target
	SynchronizeIsNeeded bool                   `json:"synchronize_is_needed"` // true, if we need to send synchronize events on releasing partitions

	OffsetPolicy OffsetPolicy `json:"offset_policy"` // specify from what topic part start message consumption
}

type OffsetPolicy string

const (
	NoOffsetPolicy      = OffsetPolicy("") // Not specified
	AtStartOffsetPolicy = OffsetPolicy("at_start")
	AtEndOffsetPolicy   = OffsetPolicy("at_end")
)

var _ model.Source = (*KafkaSource)(nil)

func (s *KafkaSource) MDBClusterID() string {
	if s.Connection != nil {
		return s.Connection.ClusterID
	}
	return ""
}

func (s *KafkaSource) WithDefaults() {
	if s.Connection == nil {
		s.Connection = &KafkaConnectionOptions{
			ClusterID:    "",
			TLS:          "",
			TLSFile:      "",
			Brokers:      nil,
			SubNetworkID: "",
		}
	}
	if s.Auth == nil {
		s.Auth = &KafkaAuth{
			Enabled:   true,
			Mechanism: "SHA-512",
			User:      "",
			Password:  "",
		}
	}
	if s.Transformer != nil && s.Transformer.CloudFunction == "" {
		s.Transformer = nil
	}
	if s.BufferSize == 0 {
		s.BufferSize = 100 * 1024 * 1024
	}
}

func (KafkaSource) IsSource() {
}

func (s *KafkaSource) GetProviderType() abstract.ProviderType {
	return ProviderType
}

func (s *KafkaSource) Validate() error {
	if s.ParserConfig != nil {
		parserConfigStruct, err := parsers.ParserConfigMapToStruct(s.ParserConfig)
		if err != nil {
			return xerrors.Errorf("unable to create new parser config, err: %w", err)
		}
		return parserConfigStruct.Validate()
	}
	return nil
}

func (s *KafkaSource) IsAppendOnly() bool {
	if s.ParserConfig == nil {
		return false
	} else {
		parserConfigStruct, _ := parsers.ParserConfigMapToStruct(s.ParserConfig)
		if parserConfigStruct == nil {
			return false
		}
		return parserConfigStruct.IsAppendOnly()
	}
}

func (s *KafkaSource) IsDefaultMirror() bool {
	return s.ParserConfig == nil
}

func (s *KafkaSource) Parser() map[string]interface{} {
	return s.ParserConfig
}

var _ model.HostResolver = (*KafkaSource)(nil)

func (s *KafkaSource) HostsNames() ([]string, error) {
	if s.Connection != nil && s.Connection.ClusterID != "" {
		return nil, nil
	}
	return ResolveOnPremBrokers(s.Connection, s.Auth)
}
