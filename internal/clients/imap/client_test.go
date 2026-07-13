package imap

import (
	"math"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/stretchr/testify/require"
)

func TestEmailCursorUIDBounds(t *testing.T) {
	zero := uint32(0)
	two := uint32(2)
	maxUid := uint32(math.MaxUint32)

	one := uint32(1)

	tests := []struct {
		name          string
		afterUid      *uint32
		beforeUid     *uint32
		wantStart     uint32
		wantStop      uint32
		wantHasCursor bool
		wantOk        bool
	}{
		{
			name:          "after_uid zero starts from the oldest message",
			afterUid:      &zero,
			wantStart:     1,
			wantStop:      0,
			wantHasCursor: true,
			wantOk:        true,
		},
		{
			name:          "after_uid at max cannot be exceeded",
			afterUid:      &maxUid,
			wantHasCursor: true,
			wantOk:        false,
		},
		{
			name:          "before_uid one leaves an empty window",
			beforeUid:     &one,
			wantHasCursor: true,
			wantOk:        false,
		},
		{
			name:          "before_uid zero leaves an empty window",
			beforeUid:     &zero,
			wantHasCursor: true,
			wantOk:        false,
		},
		{
			name:          "before_uid two bounds a single-uid window",
			beforeUid:     &two,
			wantStart:     1,
			wantStop:      1,
			wantHasCursor: true,
			wantOk:        true,
		},
		{
			name:          "neither set falls back to sequence-number path",
			wantHasCursor: false,
			wantOk:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, stop, hasCursor, ok := emailCursorUIDBounds(tt.afterUid, tt.beforeUid)

			require.Equal(t, tt.wantStart, start)
			require.Equal(t, tt.wantStop, stop)
			require.Equal(t, tt.wantHasCursor, hasCursor)
			require.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestTrimEmailPage(t *testing.T) {
	page := []domain.EmailMeta{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
	}

	t.Run("fewer than limit is unchanged", func(t *testing.T) {
		got := trimEmailPage(page, 5, false)

		require.Equal(t, page, got)
	})

	t.Run("over limit keeps the front when not keeping the newest end", func(t *testing.T) {
		got := trimEmailPage(page, 2, false)

		want := []domain.EmailMeta{{Id: "1"}, {Id: "2"}}
		require.Equal(t, want, got)
	})

	t.Run("over limit keeps the back when keeping the newest end", func(t *testing.T) {
		got := trimEmailPage(page, 2, true)

		want := []domain.EmailMeta{{Id: "2"}, {Id: "3"}}
		require.Equal(t, want, got)
	})

	t.Run("non-positive limit is unchanged", func(t *testing.T) {
		got := trimEmailPage(page, 0, false)

		require.Equal(t, page, got)
	})
}

func TestParseUid(t *testing.T) {
	t.Run("valid digits", func(t *testing.T) {
		uid, err := parseUid("42")

		require.NoError(t, err)
		require.Equal(t, uint32(42), uid)
	})

	t.Run("non-numeric string is rejected", func(t *testing.T) {
		_, err := parseUid("not-a-uid")

		require.ErrorIs(t, err, user_errors.InvalidEmailId)
	})
}
