#include<stddef.h>
#include<stdlib.h>
#include<stdio.h>
#include<assert.h>
#include<stdbool.h>

typedef int content;

typedef struct node {
    content content;
    struct node* next;
    struct node* prev;
} node;

void append_to(node* parent, content val) {
    node* child = (node*) malloc(sizeof(node));
    child->content = val;
    child->prev = parent;

    if (parent->next == NULL){
        parent->next = child;
    //insert child behind parent append what was previously appended to parent
    } else { 
        node* parent_next = parent->next;
        parent->next = child;
        child->next = parent_next;
        parent_next->prev = child;
    }
}

bool append_nodes(node* parent, node* to_append) {
    //make sure parent does not have next first
    if (parent->next != NULL) {
        return false;
    } else {
        parent->next = to_append;
        return true;
    }
}

bool delete_node(node* current) {
    node* parent = current->prev;
    node* child = current->next;

    parent->next = NULL;

    if (append_nodes(parent, child)) {
        free(current);
        return true;
    } else {
        return false;
    } 
}

/*
void delete(node* node, content to_delete) {
    //Deletes all nodes with content == to_delete
    //Delete forward
    //

    node* current = node;

    while (current->next != NULL) {
        next = current->next;

        if (current->content == to_delete) {
            node* prev = current->prev;
            prev->next = next;
            free(current)
        }

        current = next;
    }



    while (current->next != NULL) {
        next = current->next;

        if (current->content == to_delete) {
            node* prev = current->prev;
            prev->next = next;
            free(current)
        }

        current = next;
    }
}
*/

int main() {
    node* head = NULL;
    head = (node*) malloc(sizeof(node));

    if (head == NULL) {
        printf("Head still NULL");
        return 1;
    }

    head->content = 2;
    assert(head->content == 2);

    append_to(head, 4);
    assert(head->next->content == 4);

    append_to(head, 3);
    assert(head->next->content == 3);
    assert(head->next->next->content == 4);

    if (delete_node(head->next)) {
        assert(head->next->content == 4);
    };


    return 0;
}
