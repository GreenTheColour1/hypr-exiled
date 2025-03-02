{
  description = "Hypr Exiled development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
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
        fs = pkgs.lib.fileset;
      in
      {
        packages = pkgs.buildGoModule {
          pname = "hypr-exiled";
          version = "0.4.1";
          src = fs.toSource {
            root = ./.;
            fileset = fs.unions [
              ./cmd/hypr-exiled
              ./internal
              ./pkg
              ./go.mod
              ./go.sum
            ];
          };
          subPackages = [ "cmd/hypr-exiled" ];
          nativeBuildInputs = with pkgs; [
            pkg-config
          ];
          buildInputs = with pkgs; [
            alsa-lib
            xorg.libX11.dev
            xorg.libXi
            xorg.libxcb
            xorg.libXfixes
            xorg.libXext
            xorg.libXtst
          ];
          runtimeDeps = with pkgs; [
            rofi
          ];
          proxyVendor = true;
          vendorHash = "sha256-mG+aTsYvJ3Hu/CQW5Q80UV4BUDicVyzA25MC3b3pOhg=";

          meta = {
            description = "A Hyprland-focused trade manager for Path of Exile 1 and 2 with Rofi and IPC commands";
            homepage = "https://github.com/GreenTheColour1/hypr-exiled";
            platforms = pkgs.lib.platforms.linux;
          };
        };
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            # Go
            go
            gopls
            gotools
            go-tools

            xorg.libX11.dev
            xorg.libXi
            xorg.libxcb
            xorg.libXfixes
            xorg.libXext
            xorg.libXtst

            gcc

            rofi
            pkg-config
            alsa-lib
          ];

          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [ pkgs.xorg.libX11 ];

          shellHook = ''
            echo "Hypr Exiled development environment"
            echo "Ready to build with: go build -o hypr-exiled ./cmd/hypr-exiled"
          '';
        };
      }
    )
    // {
      targetSystems = [
        "aarch64-linux"
        "x86_64-linux"
      ];
    };
}
