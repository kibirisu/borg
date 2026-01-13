{
  description = "Gostack Development Environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config = {
              allowUnfree = true;
            };
          };
        in
          f pkgs
      );
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          name = "gostack-env";

          buildInputs = with pkgs; [
            pnpm
            nodejs_22
            yarn
            go
            gopls
            postgresql
            sqlc
            docker
            rootlesskit
            docker-compose
            eslint
            vtsls
            biome
            pre-commit
            golangci-lint
          ];

          shellHook = ''
            export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock
            dockerd-rootless &
            DOCKER_PID=$!
            trap "kill $DOCKER_PID 2>/dev/null" EXIT
          '';
        };
      });
    };
}
