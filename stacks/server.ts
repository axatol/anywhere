import * as k8s from "cdk8s";
import * as kplus from "cdk8s-plus-24";

class Chart extends k8s.Chart {
  generateObjectName(apiObject: k8s.ApiObject): string {
    const name = k8s.Names.toDnsLabel(apiObject, { includeHash: false });
    return name.replace(RegExp(`^${this.node.id}-`), "");
  }
}

const app = new k8s.App({
  outdir: "./dist",
  // yamlOutputType: k8s.YamlOutputType.FILE_PER_APP,
});

const chart = new Chart(app, "anywhere-server-chart", {
  namespace: "anywhere",
});

new kplus.Namespace(chart, "namespace", {
  metadata: { name: chart.namespace },
});

const deployment = new kplus.Deployment(chart, "anywhere-deploy", {
  replicas: 1,
});

const server = deployment.addContainer({
  image: "public.ecr.aws/axatol/anywhere:latest",
  name: "server",
  resources: {},
  ports: [{ number: 8042 }],
  envVariables: { CONFIG_FILEPATH: kplus.EnvValue.fromValue("/config.yml") },
  volumeMounts: [
    {
      path: "/config.yml",
      volume: kplus.Volume.fromSecret(
        chart,
        "config-vol",
        kplus.Secret.fromSecretName(chart, "config", "config")
      ),
    },
  ],
});

// const cache = deployment.addContainer({
//   image: "redis:6",
//   name: "cache",
//   ports: [{ name: "cache", number: 6379 }],
// });

const service = new kplus.Service(chart, "server-svc", {
  selector: deployment.toPodSelector(),
  ports: server.ports.map((port) => ({ port: port.number })),
});

new kplus.Service(chart, "minio-svc", {
  type: kplus.ServiceType.EXTERNAL_NAME,
  externalName: "deployment-service.minio.svc.cluster.local",
});

new kplus.Service(chart, "mongo-svc", {
  type: kplus.ServiceType.EXTERNAL_NAME,
  externalName: "deployment-service.mongo.svc.cluster.local",
});

new kplus.Ingress(chart, "server-ing").addHostRule(
  "anywhere.k8s.axatol.xyz",
  "/",
  kplus.IngressBackend.fromService(service, { port: 8042 })
);

app.synth();
