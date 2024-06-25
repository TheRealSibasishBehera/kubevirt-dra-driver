package util

import (
	"fmt"
	"github.com/google/uuid"
	"hash/fnv"
	"strings"
)

func GenerateUUIDFromSeed(seed string) string {
	h := fnv.New64()
	h.Write([]byte(seed))
	u := uuid.NewHash(h, uuid.NameSpaceDNS, []byte(seed), 0)
	return u.String()
}

func ResourceNameToEnvVar(prefix string, resourceName string) string {
	varName := strings.ToUpper(resourceName)
	varName = strings.Replace(varName, "/", "_", -1)
	varName = strings.Replace(varName, ".", "_", -1)
	return fmt.Sprintf("%s_%s", prefix, varName)
}
