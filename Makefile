PLATFORMS := darwin/arm64 windows/amd64 linux/amd64 linux/arm64

ZIPS := $(foreach p,$(PLATFORMS),target/$(subst /,_,$(p)).zip)

.PHONY: all package clean $(PLATFORMS)

all: $(PLATFORMS) $(ZIPS)

$(PLATFORMS):
	docker build --platform $@ --target export --output target/$@ .

package: $(ZIPS)

define ZIP_template
target/$(subst /,_,$1).zip: target/$1/systembolaget
	cd target/$1 && zip -q ../../$(subst /,_,$1).zip systembolaget
endef

$(foreach p,$(PLATFORMS),$(eval $(call ZIP_template,$(p))))

clean:
	rm -rf target 2>/dev/null || true
