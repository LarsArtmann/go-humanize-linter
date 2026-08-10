{
  description = "go-humanize-linter — AST linter that detects hand-rolled reimplementations of go-humanize";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";

    go-nix-helpers = {
      url = "git+ssh://git@github.com/LarsArtmann/go-nix-helpers?ref=master";
      flake = false;
    };

    go-finding = {
      url = "git+ssh://git@github.com/LarsArtmann/go-finding?ref=master";
      flake = false;
    };

    go-linter-sdk = {
      url = "git+ssh://git@github.com/LarsArtmann/go-linter-sdk?ref=master";
      flake = false;
    };

    go-error-family = {
      url = "git+ssh://git@github.com/LarsArtmann/go-error-family?ref=master";
      flake = false;
    };

    go-gogenfilter = {
      url = "git+ssh://git@github.com/LarsArtmann/gogenfilter?ref=master";
      flake = false;
    };
  };

  outputs =
    inputs@{
      self,
      nixpkgs,
      flake-parts,
      treefmt-nix,
      systems,
      ...
    }:
    let
      inherit (nixpkgs) lib;
    in
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          pkgs,
          ...
        }:
        let
          goPkg = pkgs.go_1_26;
          buildGoModule = pkgs.buildGoModule.override { go = goPkg; };

          version = self.shortRev or self.dirtyShortRev or "dev";

          mkPreparedSource = import (inputs.go-nix-helpers + "/mkPreparedSource.nix") {
            inherit pkgs lib goPkg;
          };

          preparedSrc = mkPreparedSource {
            name = "go-humanize-linter";
            inherit version;
            src = lib.cleanSourceWith {
              src = ./.;
              filter =
                path: type:
                let
                  base = baseNameOf path;
                  excludedDirs = [
                    ".direnv"
                    "result"
                    "reports"
                    "custom-gcl"
                  ];
                  isInExcludedDir = lib.any (d: lib.hasInfix "/${d}/" path) excludedDirs;
                  isExcludedBase = lib.elem base excludedDirs;
                in
                !(isInExcludedDir || isExcludedBase);
            };
            deps = {
              "github.com/larsartmann/go-finding" = inputs.go-finding;
              "github.com/larsartmann/go-linter-sdk" = inputs.go-linter-sdk;
              "github.com/larsartmann/go-error-family" = inputs.go-error-family;
              "github.com/LarsArtmann/gogenfilter/v3" = inputs.go-gogenfilter;
            };
          };

          mkApp = name: description: script: {
            type = "app";
            program = "${
              pkgs.writeShellApplication {
                inherit name;
                runtimeInputs = [
                  goPkg
                  pkgs.golangci-lint
                  pkgs.trash-cli
                ];
                text = script;
              }
            }/bin/${name}";
            meta = {
              inherit description;
              mainProgram = name;
              homepage = "https://github.com/larsartmann/go-humanize-linter";
              license = lib.licenses.mit;
              platforms = lib.platforms.unix;
              maintainers = [
                {
                  name = "Lars Artmann";
                  github = "LarsArtmann";
                }
              ];
            };
          };
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              goimports.enable = true;
              golines = {
                enable = true;
                maxLength = 120;
              };
              nixfmt.enable = true;
            };
          };

          devShells.default = pkgs.mkShell {
            packages = [
              goPkg
              pkgs.golangci-lint
              pkgs.gofumpt
              pkgs.golines
              pkgs.gopls
              pkgs.gotools
              pkgs.trash-cli
            ];

            env = {
              GOEXPERIMENT = "jsonv2";
              GOPRIVATE = "github.com/larsartmann/*,github.com/LarsArtmann/*";
            };

            shellHook = ''
              echo "go-humanize-linter dev shell — $(go version)"
              echo "GOEXPERIMENT=jsonv2 + GOPRIVATE=github.com/larsartmann/*,github.com/LarsArtmann/* active"
            '';
          };

          devShells.ci = pkgs.mkShellNoCC {
            packages = [
              goPkg
              pkgs.golangci-lint
            ];

            env = {
              GOEXPERIMENT = "jsonv2";
              GOPRIVATE = "github.com/larsartmann/*,github.com/LarsArtmann/*";
            };
          };

          packages = {
            default = buildGoModule {
              pname = "go-humanize-linter";
              inherit version;
              src = preparedSrc;
              vendorHash = "sha256-UZsYZuS7KK6mRGyEnBuj29eVrIsYYiP+bJoLFGAcihM=";
              proxyVendor = false;
              subPackages = [ "cmd/go-humanize-linter" ];
              # CLA tests shell out to `go build` via exec.Command, which requires
              # a writable HOME and can't run during buildGoModule's doCheck
              # (Nix sandbox has no writability). The upstream CI runs tests
              # separately via `nix run .#test`.
              doCheck = false;
              env.GOWORK = "off";
              preBuild = ''
                export GOEXPERIMENT=jsonv2
              '';
              ldflags = [
                "-s"
                "-w"
                "-X main.version=${version}"
              ];
              meta = {
                description = "AST linter that detects hand-rolled reimplementations of go-humanize";
                homepage = "https://github.com/larsartmann/go-humanize-linter";
                license = lib.licenses.mit;
                platforms = lib.platforms.unix;
                maintainers = [
                  {
                    name = "Lars Artmann";
                    github = "LarsArtmann";
                  }
                ];
                mainProgram = "go-humanize-linter";
              };
            };
          };

          apps = {
            test = mkApp "test" "Run all tests" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              go test ./... -count=1 "$@"
            '';

            test-race = mkApp "test-race" "Run all tests with race detector" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              go test ./... -race -count=1 "$@"
            '';

            bench = mkApp "bench" "Run benchmarks" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              go test ./... -bench=. -benchmem "$@"
            '';

            build = mkApp "build" "Build all packages and CLI" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              go build ./...
              go build -o go-humanize-linter ./cmd/go-humanize-linter/
            '';

            vet = mkApp "vet" "Run go vet" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              go vet ./...
            '';

            lint = mkApp "lint" "Run golangci-lint" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              # golangci-lint emits "Found unknown linters in //nolint
              # directives" for our //nolint:gohumanize directives because
              # the gohumanize linter is a module plugin (not a built-in).
              # To run WITH the plugin, build a custom binary:
              #   golangci-lint custom && ./custom-gcl run ./...
              # See .custom-gcl.yml for the module plugin build config.
              #
              # We capture the exit code separately from grep so lint
              # findings (exit 1) and real errors (exit >1) propagate.
              output=$(golangci-lint run ./... 2>&1)
              code=$?
              echo "$output" | grep -v 'Found unknown linters in //nolint directives'
              exit $code
            '';

            coverage = mkApp "coverage" "Run tests with coverage report" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              mkdir -p reports
              go test ./... -coverprofile=reports/coverage.out -covermode=atomic "$@"
              go tool cover -func=reports/coverage.out
            '';

            custom-lint = mkApp "custom-lint" "Build custom golangci-lint with gohumanize plugin and run it" ''
              export GOEXPERIMENT=jsonv2
              export GOPRIVATE='github.com/larsartmann/*,github.com/LarsArtmann/*'
              export GONOSUMDB='github.com/larsartmann/*,github.com/LarsArtmann/*'
              golangci-lint custom
              ./custom-gcl run -c .golangci.custom.yml ./... "$@"
            '';
          };
        };
    };
}
