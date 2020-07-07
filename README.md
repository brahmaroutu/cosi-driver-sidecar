# COSI Driver Sidecar

The cosi-driver-sidecarar is a sidecar container that runs the COSI driver so that it can run in the
user namespace without interfering with the core kubernetes code


A well known socket is opened by the COSI driver sidecar so that the provisioner can talk to the driver code to invoke cosi API.
* COSI driver socket:
  * Used by provisioner to interact with the COSI driver.
  * Created by the COSI driver.
  * Exposed on a Kubernetes node via hostpath somewhere other than the Kubelet plugin registry. (typically `/var/lib/kubelet/plugins/<drivername.example.com>/cosi.sock`).

### Required arguments

* `--cosi-address`: This is the path to the COSI driver socket (defined above) inside the
  pod that the `cosi-driver-provisioner` container will use to issue CSI
  operations (e.g. `/cosi/cosi.sock`).

### Required permissions

The cosi-driver-sidecar does not interact with the Kubernetes API, so no RBAC
rules are needed.

It does, however, need to be able to mount hostPath volumes and have the file
permissions to:

* Access the COSI driver socket (typically in `/var/lib/kubelet/plugins/<drivername.example.com>/`).
  * Used by the `cosi-driver-sidecar` to fetch the driver name from the driver

### Example

Here is an example sidecar spec in the driver DaemonSet. `<drivername.example.com>` should be replaced by
the actual driver's name.

```bash
      containers:
        - name: cosi-driver-sidecar
          image: quay.io/k8scosi/cosi-driver-sidecar:v1.0.0
          args:
            - "--cosi-address=/cosi/cosi.sock"
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "echo done"]
          volumeMounts:
            - name: plugin-dir
              mountPath: /cosi
      volumes:
        - name: plugin-dir
          hostPath:
            path: /var/lib/kubelet/plugins/<drivername.example.com>/
            type: DirectoryOrCreate
```

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

* Slack channels
  * [#sig-storage](https://kubernetes.slack.com/messages/sig-storage)
* [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-storage)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
