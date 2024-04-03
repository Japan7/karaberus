{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    japan7 = {
      url = "github:Japan7/nixpkgs-japan7";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ flake-parts, self, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];

      perSystem = { inputs', pkgs, system, ... }: {
        devShells.default = with pkgs; mkShell {
          packages = [
            go
            go-tools
            gopls
            inputs'.japan7.packages.dakara_check
            pkg-config
          ];
        };
      };
    };
}
