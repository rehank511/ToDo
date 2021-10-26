import './App.css';
import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon } from "semantic-ui-react";

let endpoint = "http://localhost:8000/todo";

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      task: "",
      items: []
    };
  }

  componentDidMount() {
    this.getTask();
  }

  onChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  onSubmit = () => {
    let { task } = this.state;
    // console.log("pRINTING task", this.state.task);
    if (task) {
      axios
        .post(endpoint + "/" + task,
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        )
        .then(res => {
          this.getTask();
          this.setState({
            task: ""
          });
          console.log(res);
        });
    }
  };

  getTask = () => {
    axios.get(endpoint).then(res => {
      console.log(res);
      if (res.data) {
        this.setState({
          items: res.data.map(item => {
            return (
              <Card key={item.Tid}>
                <Card.Content>
                  <Card.Header textAlign="left">
                    <div style={{ wordWrap: "break-word" }}>{item.Text}</div>
                    <div style={{ wordWrap: "break-word" }}>{item.Time}</div>
                  </Card.Header>

                  <Card.Meta textAlign="right">
                    <Icon
                      name="undo"
                      color="yellow"
                      onClick={() => this.updateTask(item.ID)}
                    />
                    <span style={{ paddingRight: 10 }}>Update</span>
                    <Icon
                      name="delete"
                      color="red"
                      onClick={() => this.deleteTask(item.ID)}
                    />
                    <span style={{ paddingRight: 10 }}>Delete</span>
                  </Card.Meta>
                </Card.Content>
              </Card>
            );
          })
        });
      } else {
        this.setState({
          items: []
        });
      }
    });
  };

  updateTask = id => {
    axios
      .put(endpoint + "/" + id, {
        headers: {
          "Content-Type": "application/json"
        }
      })
      .then(res => {
        console.log(res);
        this.getTask();
      });
  };

  deleteTask = id => {
    axios
      .delete(endpoint + "/" + id, {
        headers: {
          "Content-Type": "application/json"
        }
      })
      .then(res => {
        console.log(res);
        this.getTask();
      });
  };

  render() {
    return (
      <div>
          <div className="row">
            <Header className="header" as="h2">
              TO DO LIST
            </Header>
          </div>
          <div className="row">
            <Form onSubmit={this.onSubmit}>
              <Input
                type="text"
                name="task"
                onChange={this.onChange}
                value={this.state.task}
                placeholder="Create Task"
              />
              {/* <Button >Create Task</Button> */}
            </Form>
          </div>
          <div className="row">
            <Card.Group>{this.state.items}</Card.Group>
          </div>
        </div>
    );
  }
}
export default App;
