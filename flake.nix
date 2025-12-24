{
  description = "LocalCloud dev shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
    mysql = pkgs.mariadb;
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = [
        pkgs.bun
        pkgs.ffmpeg
        mysql
        pkgs.go
      ];

      shellHook = ''
        echo "Dev shell activated"

        if ! pgrep mysqld > /dev/null; then
          echo "Starting MariaDB..."
          mysql.server start
        fi
      '';
    };
  };
}
