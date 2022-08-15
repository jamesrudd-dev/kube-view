<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">

<h2 align="center">Kube View</h2>

  <p align="center">
    A display of deployed images inside your kubernetes clusters and environments
    <br />
    <br />
    <!-- <a href="https://gitlab.com/jamesrudd/kube-view">View Demo</a>
    · -->
    <a href="https://gitlab.com/jamesrudd/kube-view/-/issues">Report Bug</a>
    ·
    <a href="https://gitlab.com/jamesrudd/kube-view/-/issues">Request Feature</a>
  </p>
</div>



<!-- ABOUT THE PROJECT -->
## About The Project

Kube View was built out of necessity of enabling non-developer team members to view which version (image tags) of applications are deployed across different Kubernetes clusters and environments. Although basically built to fit my own use-requirements and selfish reasoning behind the creation, i.e. to stop non-developer team members from consistently asking which version of a specific app is in which environment, Kube View should be in a state that enables usage for any cluster set-up (configured from kubectl config file).

To note, a lot of customisation should be handled via environmental variables:
* InProduction       [bool]   - set if running in production mode or not (predominatly impacts web-server conditions) [default: `false`]
* KubeConfigLocation [string] - location of kubeconfig relative to app [default: `~/.kube/config`]
* WebServerPath      [string] - path webserver will be serving on [default: `/kube-view`]
* ImageTagFilter     [string] - set this to clean up images tags to remove (for example) the AWS prefix, i.e. ###########.dkr.ecr.ap-southeast-2.amazonaws.com/, currently only one filter available
* NamespaceFilter    [string] - comma seperated list containing namespaces desired to be removed from search, note: strings.Contains used for filter

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With:

* [![Golang][Golang]][Go-url]
* [![React][React.js]][React-url]
* [![Docker][Docker]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Docker:

Release branches are built via GitLab CI and pushed to the Docker Hub Registry.

https://hub.docker.com/r/jamesrudd/kube-view


<br>

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these example steps.

### Prerequisites

* Go v1.17
* Npm v6.14
* React
* Kubernetes Cluster (MiniKube, Docker Desktop, etc.)
* Make (optional but below instructions utilise Makefile)

### Installation

1. Clone the repo
   ```sh
   git clone https://gitlab.com/jamesrudd/kube-view.git
   ```
2. Install NPM packages
   ```sh
   make react-install
   ```
3. Install Go packages
   ```sh
   make go-install
   ```
4. Update any environmental variables to match your environment in `Makefile`
5. Run application:
   ```sh
   make run
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

[![MIT License][license-shield]][license-url]

Distributed under the MIT License. See [LICENSE](https://gitlab.com/jamesrudd/kube-view/-/blob/master/LICENSE) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

James Rudd 

E: jamesrudd.dev@gmail.com

[![LinkedIn][linkedin-shield]][linkedin-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/gitlab/contributors/jamesrudd/kube-view.svg?style=for-the-badge
[contributors-url]: https://gitlab.com/jamesrudd/kube-view/-/graphs/master
[forks-shield]: https://img.shields.io/github/forks/jamesrudd/kube-view.svg?style=for-the-badge
[forks-url]: https://gitlab.com/jamesrudd/kube-view/-/forks
[stars-shield]: https://img.shields.io/github/stars/jamesrudd/kube-view.svg?style=for-the-badge
[stars-url]: https://gitlab.com/jamesrudd/kube-view/-/starrers
[issues-shield]: https://img.shields.io/github/issues/jamesrudd/kube-view.svg?style=for-the-badge
[issues-url]: https://gitlab.com/jamesrudd/kube-view/-/issues
[license-shield]: https://img.shields.io/github/license/Ileriayo/markdown-badges?style=for-the-badge
[license-url]: https://gitlab.com/jamesrudd/kube-view/-/blob/master/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/jamesrudd15
[product-screenshot]: images/screenshot.png
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Golang]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Docker]: https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://hub.docker.com/r/jamesrudd/kube-view
