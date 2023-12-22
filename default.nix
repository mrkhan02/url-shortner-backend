# default.nix
{ lib, buildInputs, fetchDockerCompose, stdenv }:

stdenv.mkDerivation rec {
  name = "my-app";
  version = "1.0";

  buildInputs = [ fetchDockerCompose ];

  src = ./.;

  buildPhase = ''
    cp -r $src $out
    mkdir -p $out/bin
    ln -s $fetchDockerCompose/bin/docker-compose $out/bin/docker-compose
  '';

  shellHook = ''
    export COMPOSE_FILE=${src}/docker-compose.yml
  '';

  meta = with lib; {
    description = "My App";
    license = licenses.mit;
  };
}
