{
  description = "gosect - A tool to replace sections in files";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gosect = pkgs.buildGoModule {
          pname = "gosect";
          # x-release-please-start-version
          version = "0.2.0";
          # x-release-please-end
          src = ./.;

          vendorHash = "sha256-nlaO32vKmi3QVp9rZ8UCn5LIfBhLlkkiYMvuRVRK+BQ=";

          meta = with pkgs.lib; {
            description = "A tool to replace sections in files written in Go";
            homepage = "https://github.com/badele/gosect";
            license = licenses.mit;
            maintainers = [ ];
          };
        };
      in
      {
        packages = {
          default = gosect;
          gosect = gosect;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development
            go
            gopls
            gotools # goimports, godoc, etc.
            go-tools # staticcheck, etc.

            # Build tools
            just

            # Pre-commit hooks
            pre-commit

            # Docker linting
            hadolint

            gosect
          ];

          shellHook = ''
            echo "🚀 gosect development environment"
            echo "Go version: $(go version)"
            echo ""
            just
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/gosect";
        };
      }
    );
}
