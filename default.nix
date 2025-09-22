let pkgs = import <nixpkgs> { };
in pkgs.buildGoModule rec {
  pname = "releasefeed";
  version = "0.0.7";
  src = pkgs.lib.cleanSource ./.;
  vendorHash = "sha256-wroOJEkMNJpKf9OH16a4RqY8NThy6sOTt+TrdRtjpr8=";

  # network access is required for tests
  doCheck = false;
}
