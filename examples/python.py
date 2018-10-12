import os
import pwd
import getopt

# urllib with support for Python 2 and 3
try:map
    from urllib.request import urlopen, Request
    from urllib.error import HTTPError, URLError
except ImportError:
    from urllib2 import urlopen, Request, HTTPError, URLError

def parse_args(args, short_opts, long_opts):
    self.
    (parsed_opts, args) = getopt.gnu_getopt(args, short_opts, long_opts)
    opts = {}
    for (opt, value) in parsed_opts:
        opt_no_dashes = opt.lstrip('-')
        opts[opt_no_dashes] = value
    return (opts, args)

def get_user_name():
    return pwd.getpwuid(os.getuid()).pw_name

# TODO: Find a better place for these methods
def get_docker_image_info(service_params):
    docker_registry = None
    docker_image = None
    docker_tag = None

    # FIXME: Maybe this check should be done directly in definitions module
    for env, env_data in service_params.items():
        for deployment_name, deployment_params in env_data.items():
            new_docker_image = deployment_params['docker.image']
            new_docker_tag = deployment_params['docker.tag']
            if docker_registry == None:
                docker_registry = deployment_params['docker.registry']
            if docker_image == None:
                docker_image = new_docker_image
            if docker_tag == None:
                docker_tag = new_docker_tag
            if new_docker_image != docker_image or new_docker_tag != docker_tag:
                print("ERROR: Docker image/tag for {0} differs on different deployments. It should be defined on /parameters/common".format(deployment_params['service_name']))

    return (docker_registry, docker_image, docker_tag)

def check_docker_image(docker_registry, docker_image, docker_tag):
    docker_image_text = "{0}/{1}:{2}".format(docker_registry, docker_image, docker_tag)

    # Ensure the docker image exists in the repository
    docker_url = 'https://{domain}/v2/{path}/manifests/{tag}'.format(domain=docker_registry, path=docker_image, tag=docker_tag)
    headers = {"Accept": "application/vnd.docker.distribution.manifest.v2+json"}
    req = Request(docker_url, headers=headers)
    try:
        response = urlopen(req, timeout=10)
    except HTTPError as e:
        print("ERROR: Couldn't find docker image {0} (error {1})".format(docker_image_text, e.code))
        return False
    except URLError:
        print("ERROR: Couldn't connect to the docker registry to check for image {0}".format(docker_image_text))
        return False
    return True
