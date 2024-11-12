{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix

}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  packages = [
    goEnv
    gomod2nix
    pkgs.xorg.libX11
    pkgs.gcc
    pkgs.glibc
    pkgs.ffmpeg
  ];

  shellHook = ''
    export CGO_ENABLED=1
    export CGO_CFLAGS="-I${pkgs.xorg.libX11.dev}/include"
    export CGO_LDFLAGS="-L${pkgs.xorg.libX11}/lib -lX11"
    export LD_LIBRARY_PATH="${pkgs.xorg.libX11}/lib:${pkgs.gcc.cc.lib}/lib:$LD_LIBRARY_PATH"
    export LIBRARY_PATH="${pkgs.xorg.libX11}/lib:${pkgs.gcc.cc.lib}/lib:$LIBRARY_PATH"
    
    # Set the rpath for the Go binary
    export GO_LDFLAGS="-r ${pkgs.xorg.libX11}/lib -r ${pkgs.gcc.cc.lib}/lib -r ./libs"
  '';

}
