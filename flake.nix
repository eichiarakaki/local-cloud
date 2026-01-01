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
        export MYSQL_HOME=$PWD/.mysql
        export MYSQL_UNIX_PORT=$MYSQL_HOME/mysql.sock
        export MYSQL_TCP_PORT=3307

        mkdir -p $MYSQL_HOME

        if [ ! -d "$MYSQL_HOME/data/mysql" ]; then
          echo "Initializing MariaDB datadir..."
          mariadb-install-db --datadir=$MYSQL_HOME/data
        fi

        if ! pgrep -f "$MYSQL_HOME/mysql.sock" > /dev/null; then
          echo "Starting MariaDB (dev)..."
          mariadbd \
            --datadir=$MYSQL_HOME/data \
            --socket=$MYSQL_HOME/mysql.sock \
            --port=3307 \
            --skip-networking=0 \
            --pid-file=$MYSQL_HOME/mysqld.pid \
            --log-error=$MYSQL_HOME/error.log &
        fi

        echo "MariaDB running on socket $MYSQL_UNIX_PORT"
      '';
    };
  };
}
