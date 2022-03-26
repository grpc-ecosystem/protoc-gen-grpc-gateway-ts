package registry

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// CommentInfo organizes SourceCodeInfo proto into a nested structure for easier access.
//
// Each CommentInfo represents the comments in the context of a prefix of SourceCodeInfo.Location.path.
// See the definition of the path at: https://github.com/protocolbuffers/protobuf/blob/236d04c5b4431bd47bc88402e2a9332ba52847d6/src/google/protobuf/descriptor.proto#L767
type CommentInfo struct {
	// text is the text for this source context.
	text string
	// SubComments is a map from the next path element to the corresponding CommentInfo reprenting it.
	subComments map[int]*CommentInfo
}

// AddLocation adds a Location proto to the root CommentInfo.
//
// DO NOT call it on a child CommentInfo
func (c *CommentInfo) AddLocation(l *descriptorpb.SourceCodeInfo_Location) {
	current := c
	for _, path := range l.GetPath() {
		child := current.GetSubComment(int(path))
		if child == nil {
			child = &CommentInfo{}
			current.putSubComment(path, child)
		}
		current = child
	}
	commentLines := []string{}
	for _, comment := range l.GetLeadingDetachedComments() {
		commentLines = appendComment(commentLines, comment)
	}
	commentLines = appendComment(commentLines, l.GetLeadingComments())
	commentLines = appendComment(commentLines, l.GetTrailingComments())
	if len(commentLines) > 0 {
		current.text = "/**\n" + strings.Join(commentLines, "\n") + "\n **/"
	}
}

func (c *CommentInfo) GetSubComment(path ...int) *CommentInfo {
	current := c
	for _, p := range path {
		if current == nil {
			return nil
		}
		if current.subComments == nil {
			return nil
		}
		subComment, ok := current.subComments[p]
		if !ok {
			return nil
		}
		current = subComment
	}
	return current
}

func (c *CommentInfo) GetText() string {
	if c == nil {
		return ""
	}
	return c.text
}

func (c *CommentInfo) putSubComment(path int32, child *CommentInfo) {
	if c.subComments == nil {
		c.subComments = make(map[int]*CommentInfo)
	}
	c.subComments[int(path)] = child
}

func (c *CommentInfo) String() string {
	return c.string([]string{})
}

func (c *CommentInfo) string(prefix []string) string {
	ret := fmt.Sprintf("[%s]: %s\n", strings.Join(prefix, ", "), c.text)
	childPrefix := append(prefix[:], "")
	for k, v := range c.subComments {
		childPrefix[len(prefix)] = strconv.FormatInt(int64(k), 10)
		ret += v.string(childPrefix)
	}
	return ret
}

func appendComment(lines []string, comment string) []string {
	if strings.TrimSpace(comment) == "" {
		return lines
	}
	if len(lines) > 0 {
		lines = append(lines, " *")
	}
	split := strings.Split(comment, "\n")
	for len(split) > 0 && strings.TrimSpace(split[len(split)-1]) == "" {
		split = split[:len(split)-1]
	}
	for _, line := range split {
		lines = append(lines, " *"+line)
	}
	return lines
}
