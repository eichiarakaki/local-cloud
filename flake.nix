{ pkgs ? import <nixpkgs> {} }:

let
  mysql = pkgs.mariadb;
in

pkgs.mkShell {
  buildInputs = [
    pkgs.bun
    pkgs.ffmpeg
    mysql
    go
    
  ];

  shellHook = ''
    # Verificar si estamos en el entorno de desarrollo y activar MySQL
    if [ -z "$NIX_BUILD_TOP" ]; then
      echo "Activando MySQL solo en desarrollo"
      # Comando para iniciar MySQL si no está en ejecución
      if ! pgrep mysqld > /dev/null; then
        echo "Iniciando MySQL..."
        # Comando para iniciar el servicio MySQL
        mysql.server start
      fi
    fi
  '';
}
