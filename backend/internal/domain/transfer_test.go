package domain

import (
	"testing"
	"time"
)

func TestTransferIsExpired(t *testing.T) {
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	t.Run("status expired", func(t *testing.T) {
		tr := Transfer{
			Status:    TransferStatusExpired,
			ExpiresAt: now.Add(time.Hour),
		}
		if !tr.IsExpired(now) {
			t.Fatal("expected expired when status is expired")
		}
	})

	t.Run("expires_at in the past", func(t *testing.T) {
		tr := Transfer{
			Status:    TransferStatusActive,
			ExpiresAt: now.Add(-time.Second),
		}
		if !tr.IsExpired(now) {
			t.Fatal("expected expired when expires_at is past")
		}
	})

	t.Run("expires_at equal to now", func(t *testing.T) {
		tr := Transfer{
			Status:    TransferStatusActive,
			ExpiresAt: now,
		}
		if !tr.IsExpired(now) {
			t.Fatal("expected expired when expires_at == now")
		}
	})

	t.Run("still active", func(t *testing.T) {
		tr := Transfer{
			Status:    TransferStatusActive,
			ExpiresAt: now.Add(time.Minute),
		}
		if tr.IsExpired(now) {
			t.Fatal("expected not expired")
		}
	})
}

func TestTransferEffectiveStatus(t *testing.T) {
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		tr   Transfer
		want string
	}{
		{
			name: "deleted wins over expiry",
			tr: Transfer{
				Status:    TransferStatusDeleted,
				ExpiresAt: now.Add(-time.Hour),
			},
			want: TransferStatusDeleted,
		},
		{
			name: "blocked wins over time expiry",
			tr: Transfer{
				Status:    TransferStatusBlocked,
				ExpiresAt: now.Add(-time.Hour),
			},
			want: TransferStatusBlocked,
		},
		{
			name: "time expiry overrides active",
			tr: Transfer{
				Status:    TransferStatusActive,
				ExpiresAt: now.Add(-time.Minute),
			},
			want: TransferStatusExpired,
		},
		{
			name: "active when not expired",
			tr: Transfer{
				Status:    TransferStatusActive,
				ExpiresAt: now.Add(time.Hour),
			},
			want: TransferStatusActive,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.tr.EffectiveStatus(now)
			if got != tc.want {
				t.Fatalf("EffectiveStatus = %q, want %q", got, tc.want)
			}
		})
	}
}
