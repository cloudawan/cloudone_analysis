go build
mv kubernetes_management_analysis docker/kubernetes_management_analysis/
find ! -wholename './docker/*' ! -wholename './docker' ! -wholename '.' -exec rm -rf {} +
mv docker/version version
mv docker/environment environment