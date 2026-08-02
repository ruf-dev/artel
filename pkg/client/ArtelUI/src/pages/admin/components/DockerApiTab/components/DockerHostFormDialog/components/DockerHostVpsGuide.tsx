import cls
from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/components/DockerHostVpsGuide.module.css"

export default function DockerHostVpsGuide() {
    return (
        <details className={cls.DockerHostVpsGuideContainer}>
            <summary className={cls.Summary}>How do I set up a public Docker API on a VPS?</summary>
            <p>1. Install Docker on the VPS if it isn't already.</p>
            <p>
                2. Generate a CA plus a server and client certificate/key pair — mutual TLS is required, since an
                unauthenticated Docker API is unauthenticated root access to the host. See Docker&apos;s official
                guide: <a href="https://docs.docker.com/engine/security/protect-access/" target="_blank"
                    rel="noopener noreferrer">docs.docker.com/engine/security/protect-access</a>.
            </p>
            <p>3. Start (or restart) the daemon with TLS enabled and listening on the public interface:</p>
            <pre className={cls.Command}>
                <code>
                    {"dockerd --tlsverify \\\n"}
                    {"  --tlscacert=ca.pem --tlscert=server-cert.pem --tlskey=server-key.pem \\\n"}
                    {"  -H=0.0.0.0:2376"}
                </code>
            </pre>
            <p>4. In the VPS firewall, open port 2376 — ideally restricted to Artel&apos;s server IP only.</p>
            <p>
                5. Above, set URL to <code className={cls.Inline}>tcp://&lt;vps-ip&gt;:2376</code>, then paste
                {" "}<code className={cls.Inline}>ca.pem</code>, <code className={cls.Inline}>client-cert.pem</code>,
                {" "}and <code className={cls.Inline}>client-key.pem</code> into the CA certificate, Client
                certificate, and Client key fields below.
            </p>
        </details>
    )
}
