{
  description = "Adventures on running the Svelve compiler in Go";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        buildDeps = with pkgs; [
          git
          go_1_20
          gnumake
          esbuild
          nodejs_18
        ];
        devDeps = with pkgs;
          buildDeps ++ [
            wasmtime
          ];
      in
      { devShell = pkgs.mkShell { buildInputs = devDeps; }; }
    );

}
