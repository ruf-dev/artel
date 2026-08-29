package simplechat

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
)

func TestDeriveTitle_EmptyText(t *testing.T) {
	require.Equal(t, "", deriveTitle(""))
}

func TestDeriveTitle_BlankText(t *testing.T) {
	require.Equal(t, "", deriveTitle("   \n\t  "))
}

func TestDeriveTitle_ShortSingleLine(t *testing.T) {
	require.Equal(t, "How do I set up CouchDB?", deriveTitle("How do I set up CouchDB?"))
}

func TestDeriveTitle_TrimsSurroundingWhitespace(t *testing.T) {
	require.Equal(t, "Trim me", deriveTitle("   Trim me   "))
}

func TestDeriveTitle_MultiLineKeepsOnlyFirstLine(t *testing.T) {
	text := "First line is the topic\nSecond line is detail\nThird line too"

	require.Equal(t, "First line is the topic", deriveTitle(text))
}

// TestDeriveTitle_TruncatesLongTextToMaxTitleLen guards the exact cap: maxTitleLen runes kept,
// plus the ellipsis marker, for text past the limit.
func TestDeriveTitle_TruncatesLongTextToMaxTitleLen(t *testing.T) {
	text := strings.Repeat("a", maxTitleLen+40)

	title := deriveTitle(text)

	expected := strings.Repeat("a", maxTitleLen) + "…"
	require.Equal(t, expected, title)
	require.Equal(t, maxTitleLen+1, len([]rune(title)))
}

// TestMaybeSetTitle_AlreadyTitledIsNoOp confirms the chat.Title guard short-circuits before
// maybeSetTitle ever reaches mutateChatFile — every dependency mutateChatFile would need
// (vaults, couch accounts/instances, a live CouchDB) is left nil on s, so reaching that code
// would panic on the nil vaults repo. A zero-value Service surviving the call unharmed is the
// proof the guard fired first.
func TestMaybeSetTitle_AlreadyTitledIsNoOp(t *testing.T) {
	existingTitle := "Already named"
	chat := domain.SimpleChat{
		Uuid:  uuid.New(),
		Title: &existingTitle,
	}

	s := &Service{}

	require.NotPanics(t, func() {
		s.maybeSetTitle(context.Background(), chat, "does this ever get named twice?")
	})
}

// TestMaybeSetTitle_EmptyDerivedTitleIsNoOp confirms a message that derives to an empty title
// (blank/whitespace-only) also short-circuits before touching mutateChatFile.
func TestMaybeSetTitle_EmptyDerivedTitleIsNoOp(t *testing.T) {
	chat := domain.SimpleChat{
		Uuid: uuid.New(),
	}

	s := &Service{}

	require.NotPanics(t, func() {
		s.maybeSetTitle(context.Background(), chat, "   ")
	})
}
