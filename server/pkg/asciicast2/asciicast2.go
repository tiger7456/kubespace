/*




Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package asciicast2

import (
	"bytes"
	"encoding/json"
	"kubespace/server/pkg/utils"
)

type CastV2Header struct {
	Version      uint               `json:"version"`
	Width        int                `json:"width"`
	Height       int                `json:"height"`
	Timestamp    int64              `json:"timestamp,omitempty"`
	Duration     float64            `json:"duration,omitempty"`
	Title        string             `json:"title,omitempty"`
	Command      string             `json:"command,omitempty"`
	Env          *map[string]string `json:"env,omitempty"`
	outputStream *json.Encoder
}

func NewCastV2(meta CastV2Header, stream *bytes.Buffer) (*CastV2Header, *bytes.Buffer) {
	var c CastV2Header
	c.Version = 2
	//c.Width = meta.Width
	//c.Height = meta.Height
	// 固定宽高用于前端展示
	c.Width = meta.Width
	c.Height = meta.Height
	c.Title = meta.Title
	c.Timestamp = meta.Timestamp
	c.Duration = c.Duration
	c.Env = meta.Env
	c.outputStream = json.NewEncoder(stream)
	c.outputStream.Encode(c)
	return &c, stream
}

func (c *CastV2Header) Record(t float64, data []byte, event string) {
	out := make([]interface{}, 3)
	//timeNow := time.Since(t).Seconds()
	out[0] = t
	out[1] = event // i：input；o：output
	out[2] = utils.Bytes2Str(data)
	c.Duration = t
	c.outputStream.Encode(out)
}
