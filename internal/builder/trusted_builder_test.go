package builder_test

import (
	"testing"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	bldr "github.com/buildpacks/pack/internal/builder"
	"github.com/buildpacks/pack/internal/config"

	h "github.com/buildpacks/pack/testhelpers"
)

func TestTrustedBuilder(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "Trusted Builder", trustedBuilder, spec.Parallel(), spec.Report(report.Terminal{}))
}

func trustedBuilder(t *testing.T, when spec.G, it spec.S) {
	when("IsKnownTrustedBuilder", func() {
		it("matches exactly", func() {
			h.AssertTrue(t, bldr.IsKnownTrustedBuilder("paketobuildpacks/builder-jammy-base"))
			h.AssertFalse(t, bldr.IsKnownTrustedBuilder("paketobuildpacks/builder-jammy-base:latest"))
			h.AssertFalse(t, bldr.IsKnownTrustedBuilder("paketobuildpacks/builder-jammy-base:1.2.3"))
			h.AssertFalse(t, bldr.IsKnownTrustedBuilder("my/private/builder"))
		})
	})

	when("IsTrustedBuilder", func() {
		it("matches partially", func() {
			cfg := config.Config{
				TrustedBuilders: []config.TrustedBuilder{
					{
						Name: "my/trusted/builder-jammy",
					},
				},
			}
			builders := []string{
				"my/trusted/builder-jammy",
				"my/trusted/builder-jammy:latest",
				"my/trusted/builder-jammy:1.2.3",
			}

			for _, builder := range builders {
				isTrusted, err := bldr.IsTrustedBuilder(cfg, builder)
				h.AssertNil(t, err)
				h.AssertTrue(t, isTrusted)
			}
			isTrusted, err := bldr.IsTrustedBuilder(cfg, "my/private/builder")
			h.AssertNil(t, err)
			h.AssertFalse(t, isTrusted)

			isTrusted, err = bldr.IsTrustedBuilder(cfg, "my/trusted/builder-jammy-base")
			h.AssertNil(t, err)
			h.AssertFalse(t, isTrusted)
		})
	})
}
